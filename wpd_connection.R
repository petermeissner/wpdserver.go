wpds_connection <- function(){
  
  # read in credentials and split in to key/values
  split <- 
    strsplit(
      strsplit(
        readLines(".db_credentials"), 
        " "
      )[[1]], 
      "="
    )
  
  # save key value pairs
  creds        <- vapply(split, `[`, character(1), 2)
  names(creds) <- vapply(split, `[`, character(1), 1)
  
  # open up connection
  con <- 
    dbConnect(
      PostgreSQL(), 
      user     = creds["user"], 
      password = creds["password"], 
      dbname   = creds["dbname"],
      host     = creds["host"]
    )
}
