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

### Query
```
```


### Insert
```
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
    return "measurement1"
}

```
The function ```Bucket()``` and ```Measurement()``` are requested, telling the bucket / measurement where the record should be.
 
### Tags
The field tag ```florm:"k,xxx"``` and ```florm:"v,xxx"``` indicates wehter this field if a **tag** or **field**, as well as their binded name in InfluxDB. 

> **WARNING**: Please keep in mind that InfluxDB is SCHEMALESS! So, do NOT expect those data models above can be as reliable as regular relational DB (mysql). 
> - Impossible to change schema for those ***"past"*** records. Either use them in a compatible way or just simply delete them. 
> - The value of tag columns can NOT be modified. If you really need it, please **delete** it and **insert** again.

## FluxSession


## Insert

## Update
> **WARNING:** The Florm update function follows the "Query - Set - Write" pipeline. Data writing is implemented by the ```to()``` clause of Flux language. You can select any data to be updated using flux query clauses and we will protect the update by filtering out those columns that are not in the given datamodel. However, we are not able to recall those columns that are missed by your query. So please review your update script **VERY CAREFULLY** before you run it, otherwise it will cause unpredicatable data inconsistency..

## Delete
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