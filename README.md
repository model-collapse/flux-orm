# Florm: An ```ORM``` tool for flux language in Go.

InfluxDB is a widely used time-series database for IoT, micro-service and finance scenarios. InfluxDB 2.0 has released the Flux language to althernate InfluxQL with some powerful features such as stream-based (DAG) data aggregation, lambda experssion and more transformation functions.

```Florm``` is a ORM framework to make it easier for developers to use flux inside Go apps. Just download:
```
go get -i github.com/model-collapse/flux-orm
```
and import it in your code:
```
import "github.com/model-collapse/flux-orm/florm"
```

## New Features
- Inline flux query builder, to write flux style queries in native Go! 
- Supporting most of the Flux keywords and syntax. Keep improving.
- Same with ```GORM```, people can define data schemas inside go structs. You can use structs as data receivers of your query results.
- Partially repurposed InfluxDB from a pure time-series DB to a multi-modal DB by adapting non-series datamodel. See secion "Data Models" below for more imformation.

## Get Started
There are some initial examples for you to understand how to use Florm. 

### Insert:
```
import (
    "github.com/model-collapse/flux-orm/florm"
    "time"
    "log"
)

type Log struct {
    Series
    IsFatal bool   `florm:"k,fatal"`
    Position int   `florm:"k,position"`
    Code string    `florm:"k,code"`
    Risk float32   `florm:"v,risk"`
    Content string `florm:"v,content"`
}

func (l *Log) Bucket() string {
    return "misc1"
}

func (l *Log) Measurement() string {
    return "log"
}

ss := florm.NewFluxSession()

logs := []Log {
    {
        Series: Series{
            Time: time.Now()
        },
        IsFatal: false,
        Position: 1,
        Code: "json",
        Risk: 0.918,
        Content: "[Info] it is running",
    },
    {
        Series: Series{
            Time: time.Now()
        },
        IsFatal: true,
        Position: 2,
        Code: "c++",
        Risk: 1.97,
        Content: "[Fatal] it is NOT running",
    }
}

if err := ss.Insert(logs); err != nil {
    log.Fatal(err)
}
```

### Query:
```
type RiskAgg struct {
    Code string `florm:"k,code"`
    Risk float32 `florm:"v,risk"`
}

ss := florm.NewFluxSession()

var riskHistogram []RiskAgg
start := time.Date(2021, 7, 7, 0, 0, 0, 0, time.UTC)
stop := time.Now()
ss.Range(start, stop).Filter("(r) => (r._measurement==\"log\" and r._field==\"risk\")", "drop").GroupBy([]string{"code"}).Sum().Yield(&riskHistogram)

if err := ss.ExecuteQuery(context.Background()); err != nil {
    log.Fatal(err)
}

```

## Data Models
There are **3** kinds of data models in ```Florm```: *Series*, *DTable*, *STable*
| Model | Description |
| ----- | ----------- |
| Series | This is the traditional time-series data item in InfluxDB, each item is assigned with a primary tag **"_time"** as index. |
| STable | This is a data model for relational data items. The **"_time"** tag is always '0' in the system, and with a manually assigned primary key called **"primary"** as a tag. It is recommended to be used in storing data that is not frequently inserted, like enumeration mapping. |
| DTable | This ia also a data model for relational data items. The difference with ```STable``` is that ```DTable``` has both **"_time"** tag and **"_primary"**. The tag **"_time"** is binded with the time upon insertion and **"_primary"** ID comes from *snowflake* ID generation. It is suitable for those tables with dynamic growing sizes.

### Usage 
Just embbed those models inside our custom data model with composition:
```
type Student struct {
    florm.STable
    Name string `florm:"k,name"`
    Class int `florm:"k,class"`
    Score float64 `florm:"v,score"`
}

func (s *Student) Bucket() string {
    return "sample-bucket"
}

func (s *Student) Measurement() string {
    return "grades"
}

```
The function ```Bucket()``` and ```Measurement()``` are requested, telling the bucket / measurement where the record should be.
 
### Tags
The field tag ```florm:"k,xxx"``` and ```florm:"v,xxx"``` indicates wehter this field if a **tag** or **field**, as well as their binded name in InfluxDB. 

> **WARNING**: Please keep in mind that InfluxDB is SCHEMALESS! So, do NOT expect those data models above can be as reliable as regular relational DB (mysql). 
> - Impossible to change schema for those ***"past"*** records. Either use them in a compatible way or just simply delete them. 
> - The value of tag columns can NOT be modified. If you really need it, please **delete** it and **insert** again.

## FluxSession
Each FluxSession is recommended to be used only once, if you want to use it twice, the query content in the last execution will still be there. So when you do anothe query, create a new one.
The user can also create a session with customized ```APIManager``` object in order to manage APIs yourself. 
```
var mgr florm.APIManager
...

ss := NewFluxSessionCustomAPI(mgr)
```
The default api manager can be set via ```florm.RegisterDefaultAPIManager()```.

## Insert
Inserting with Florm is mentioned in the example above. Besides an slice of structs, the insert api can also accept follow types:
```
// Single Struct
Log{...}

// Pointer to Struct
&Log{...}

// Slice of Struct
[]Log{...}

// Slice of Struct Points
[]*Log{...}

// Channel of struct 
<-chan Log 
<-chan *Log
``` 
> **WARNING:** The inserting process will not terminate util channels are closed, so be careful when use channels.

## Update
Here is an example of data updating.
```
ss := NewFluxSession()

s := Student {
    Score: -1.0,
}

// Static is a range selection clause specificly for STable models
ss.Static().Filter(`(r) => (r._measurement=="grades" and r.name=="sean")`, "drop").Update([]string{"name"}, &s)

if err := ss.ExecuteQuery(context.Background()); err != nil {
    log.Fatal(err)
}
```

> **WARNING:** The Florm update function follows the "Query - Set - Overwrite" pipeline. Data writing is implemented by the ```to()``` clause of Flux language. You can select any data to be updated using flux query clauses even with ```join```, and we will protect the update by filtering out those columns that are not in the given datamodel. However, we are not able to recall those columns that are missed by your query. So please review your update script **VERY CAREFULLY** before you run it, otherwise it will cause unpredicatable data inconsistency..

## Delete
Example of deleting some records:
```
ss := NewFluxSession()

// deleting with primary ids
if err := ss.Delete(context.Background(), &Student{}, 1, 2, 3); err != nil {
    log.Fatal(err)
} 

if err := ss.DeleteWithFilter(context.Background(), &Log{Series:Series{Start:time.Unix(0,0), Stop:time.Now()}}, "r._measurement==\"log\""); err != nil {
    log.Fatal(err)
}
```

> **Notes**: the lambda expression used by delete API is different with those used in queries. Instead of using ```(r) => (r._measurement=="??")```, please just use the condition clause ```r._measurement=="??"```. This is suggested by InfluxDB offical API.


## FAQ
### Why do we need to repurpose InfluxDB to be multi-modal? Why not directly use relational DB?
***Ans***: **Close this page** and sleep with Postgres.

### My InfluxDB cluster are configured with complicated token management, can Florm handle this.
***Ans***: You can customize the ```APIManager``` interface to assign the best api object for each given bucket.

### My query results are of newly created schema, which is consistent with those in the DB. Should I create a new datamodel?
***Ans***: You need a new struct binding the schema of the result using field tags, but it does not need to be composed with a datamodel type. 

### How about those advanced features of Flux, such as array or map definition?
***Ans***: We didn't put thos inside this version, and we need more discussions about those non-inline syntax. Feel free to drop your thoughs in the issues.

### Considering multi-model DB, the Flux language already supported call SQL-DB with native api. Can we support that in the ORM?
***Ans***: It is in the plan, will be very impressive if Florm works together with GORM.