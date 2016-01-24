
# Copyright (C) 2016 J. Robert Wyatt
# 
# This program is free software; you can redistribute it and/or
# modify it under the terms of the GNU General Public License
# as published by the Free Software Foundation; either version 2
# of the License, or (at your option) any later version.
# 
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
# 
# You should have received a copy of the GNU General Public License
# along with this program; if not, write to the Free Software
# Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.



struct Customer {
    1: string FirstName,
    2: string LastName,
    3: string Addr1,
    4: string Addr2,
    5: string City,
    6: string StateProvince,
    7: string PostalZip,
    8: string Country,
}



struct Result {
    1: bool     Ok,     # True means success, false otherwise
    2: string   Key,    # Optional. Needed for create operations
    3: string   Err,    # Optional. Error string when ok=false
}

service DataStore {
            # Update the given customer with any changed details
            Result                  CreateCustomer(1:Customer customer),
            # Get Customer by key
            Customer                GetCustomer(1:string key),
            # Get all customer keys matching details in template
            map<string,Customer>    GetAllCustomers(1:Customer template),
            # Update the given customer with any changed details
            Result                  UpdateCustomer(1:string key, 2:Customer customer),
}
