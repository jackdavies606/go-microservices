# go-microservices

##Endpoints;

###Customer Service : 8080

GetCustomers - GET - /customers

GetCustomerByName - GET - /customer/name/{name}

GetCustomerById - GET - /customer/id/{id}

AddCustomer - POST - /customer/name/{name}

RemoveCustomer - DELETE - /customer/name/{name}

###Item Service : 8081

GetItems - GET - /items

GetItemByName - GET - /item/name/{name}

GetItemById - GET - /item/id/{id}

RemoveItem - DELETE - /item/name/{name}

AddItem - POST - /item

###Order Service : 8082

GetAllCustomerOrders - GET - /orders/customer/{customerId}

GetCustomersOpenOrder - GET - /order/customer/{customerId}

AddToOrder - POST - /order/customer/{customerId}/item/{itemId}

CancelOrder - DELETE - /order/customer/{customerId}
