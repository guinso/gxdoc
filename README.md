# gxdoc
evaluate golang on building up traceable document system

## Web server configuration
1. run gxdoc 
2. a config.ini file will be generated at same directory of gxdoc executable binary
3. edit config.ini with any text editor

## REST API
### API Summary
|HTTP Method|URL|Description|
| --- | --- | --- |
| GET | /api/document/schema-infos | get all schema info |
| POST | /api/document/schema-infos | create new schema info |
| GET | /api/document/schema-infos/{schema-name} | get specific schema info |
| POST | /api/document/schema-infos/{schema-name} | update schema info |
| GET | /api/document/schemas/{schema-name}/revisions/{revision-number} | get specific schema definition by revision number |
| GET | /api/document/schemas/{schema-name} | get latest schema definition |
| POST | /api/document/schemas/{schema-name} | update schema definition to latest revision |
| GET | /api/document/schemas/{schema-name}/draft | get draft version of schema definition |
| POST | /api/document/schemas/{schema-name}/draft | update draft version of schema definition | 

### Get list of Schema Infomation
NOTE: <i><b>{dev-start-url}</b> is defined in config.ini file</i>

URL Pattern:    
```
GET /api/document/schema-infos
```
Output:
```json
{
    "statusCode": 0,
    "statusMsg": "ok",
    "response": [
        {
            "name": "pr",
            "latestRev": 1,
            "desc": "purchase requisite",
            "isActive": true,
            "hasDraft": true
        },
        {
            "name": "invoice",
            "latestRev": 3,
            "desc": "invoice....",
            "isActive": true,
            "hasDraft": false
        }
    ]
}
```

### Get Schema Information By Schema Name
URL Pattern:
```
GET /api/document/schema-infos/{schema-name}
```
Output:
```json
{
    "statusCode": 0,
    "statusMsg": "",
    "response": {
        "name": "pr",
        "latestRev": 1,
        "desc": "purchase requisite",
        "isActive": true,
        "hasDraft": true
    }
}
```

### Register New Schema Information
URL Pattern:
```
POST /api/document/schema-infos
```
Input Data (sample):
```json
{
    "name":"po",
    "desc":"purchase order"
}
```

### Update Schema Information
URL Pattern:
```
POST /api/document/schema-infos/{schema-name}
```
Input Data (sample):
```json
{
    "name":"po",
    "desc":"new PO description",
    "isActive": true
}
```

### Get Latest Schema Definition
URL Pattern:
```
GET /api/document/schemas/{schema-name}
```
Output:
```xml
<?xml version="1.0"?>
<dxdoc name="invoice" revision="3" id="">
    <dxstr name="invNo"></dxstr>
    <dxint name="totalQty" isOptional="true"></dxint>
    <dxdecimal name="price" precision="2"></dxdecimal>
</dxdoc>
```

### Get Schema Definition by Revision
URL Pattern:
```
GET /api/document/schemas/{schema-name}/revision/{revision-number}
```
Output:
```xml
<?xml version="1.0"?>
<dxdoc name="invoice" revision="2" id="">
    <dxstr name="invNo"></dxstr>
    <dxint name="totalQty" isOptional="true"></dxint>
    <dxdecimal name="price" precision="6"></dxdecimal>
</dxdoc>
```

### Update Schema Definition
NOTE: <i>newly update schema definition will register as new definition with higher revision number. previous definition will remain intact in database</i>

URL Pattern:
```
POST /api/document/schemas/{schema-name}
```
Input Data (sample):
```xml
<?xml version="1.0"?>
<dxdoc name="invoice" revision="0" id="">
    <dxstr name="invNo"></dxstr>
    <dxint name="totalQty" isOptional="true"></dxint>
    <dxdecimal name="price" precision="2"></dxdecimal>
    <dxbool name="needAudit"></dxbool>
</dxdoc>
```

### Get Schema Definition's Draft
URL Pattern:
```
GET /api/document/schemas/{schema-name}/draft
```
Output:
```xml
<?xml version="1.0"?>
<dxdoc name="invoice" revision="-1" id="">
    <dxstr name="invNo"></dxstr>
    <dxint name="totalQty" isOptional="true"></dxint>
    <dxdecimal name="price" precision="2"></dxdecimal>
</dxdoc>
```

### Update Schema Definition's Draft
NOTE: <i>newly posted schema definition will overwrite previous draft definition!</i>

URL Pattern:
 ```
POST /api/document/schemas/{schema-name}/draft
```
Input Data (sample):
```xml
<?xml version="1.0"?>
<dxdoc name="invoice" revision="0" id="">
    <dxstr name="invNo"></dxstr>
    <dxint name="totalQty" isOptional="true"></dxint>
    <dxdecimal name="price" precision="2"></dxdecimal>
</dxdoc>
```
