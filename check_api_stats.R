library(RPostgreSQL)
library(dplyr)


source("wpd_connection.R")


con <- wpds_connection()


rs <- dbSendQuery(con, "SELECT * from api_usage")
dbColumnInfo(rs)

dbReadTable(con,"api_usage") %>% tail(100)


