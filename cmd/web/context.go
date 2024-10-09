package main

type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")
const userContextKey = contextKey("user")
