library(RPostgreSQL)

source("wpd_connection.R")


con <- wpds_connection()


dbCreateTable(
  conn   = con, 
  name   = "api_usage", 
  fields = 
    data.frame(
      api      = "test",
      resource = "test",
      ts       = as.POSIXct("2019-01-01 00:00:00")
    )
)

