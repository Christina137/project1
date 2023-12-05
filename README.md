## project1
My first web app project

### ShowVideo:

https://github.com/Christina137/project1/assets/98719024/cc230b12-ecff-4515-9644-d538bd6d5bc5


### To Run this project:
I am using Oracle Cloud Infrastructure (OCI) for object storege service. You have to prepare the following information for this src/dao/OracleCloud.go file:
```go
const TenancyOCID string = "your-tenancy-ocid"
const UserOCID string = "your-user-ocid"
const Region string = "your-region"
const Fingerprint string = "your-fingerprint"
const compartmentOCID string = "your-compartment-ocid"
const PrivateKey string = "your-private-key"
const Namespace string = "your-namespace"
const BucketName string = "your-bucket-name"
const Url string = "your-bucket-url"
```
next, Set up your mysql database information in this resources/application.yml file :
```
mysql:
  url: "your-ip"
  userName: root
  passWord: "your-password"
  dbname: "your-dbname"
  port: 3306
jwt:
  secret: "your-secret"
```

#### Acknowledgment：[youthcamp@bytedance.com](https://youthcamp.bytedance.com/)

