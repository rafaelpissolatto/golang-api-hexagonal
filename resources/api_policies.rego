package golangapitemplate.authz

# By default, deny requests.
default allow = false

# Admins can do everything
allow {
	input.Token.Roles[_] == "admin"
}

# Business are allowed to update product
allow {
    input.EntityData.Type == "updateProduct"
	input.Token.Roles[_] == "business"
}

# Business are allowed to view product
allow {
    input.EntityData.Type == "viewProduct"
	input.Token.Roles[_] == "business"
}

# User are allowed to view product
allow {
    input.EntityData.Type == "viewProduct"
	input.Token.Roles[_] == "user"
}

# This business Username is allowed to create product
allow {
    input.EntityData.Type == "createProduct"
	input.Token.Roles[_] == "business"
	input.Token.Username == "business_main_id"
}

# Only owners can edit objects
allow {
    input.Token.Username == input.EntityData.Owner
}
