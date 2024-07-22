variable "database" {
  type = string
}
variable "schema" {
  type = string
}
variable "name" {
  type = string
}
variable "stage" {
  type = string
}
variable "directory_location" {
  type = string
}
variable "main_file" {
  type = string
}
variable "query_warehouse" {
  type = string
}
variable "external_access_integrations" {
  type = list(string)
}
variable "title" {
  type = string
}
variable "comment" {
  type = string
}
