package enums

const (
	BDMongo      string = "DB.Mongo"
	DBConnection string = "DB.getConnection"

	MongodbMaxDatabaseTimeOut    int64  = 30000
	MongodbMaxConnectionTimeOut  int64  = 3000
	MongodbSocketTimeout         int64  = 300000
	MongodbSocketReadTimeout     int64  = 200000
	MongodbMaxConnectionIdleTime int64  = 300000
	MongodbMinConnectionsPerHost uint64 = 0
	MongodbMaxConnectionsPerHost uint64 = 150

	MongodbDatabase        string = "data-master"
	LoginLogOTP            string = "login-log-otp"
	LoginBlackListOTP      string = "login-black-list-otp"
	LoginFailedAttemptsOTP string = "login-failed-attempts-otp"
)
