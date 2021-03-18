query, err := ioutil.ReadFile("path/to/database.sql")
if err != nil {
panic(err)
}
if _, err := db.Exec(query); err != nil {
panic(err)
}

