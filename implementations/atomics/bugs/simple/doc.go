/*
Package simple provides the intentionally incorrect atomic balance that
only tracks the raw value. It is designed to highlight how missing CAS
protection fails under contention.
*/
package simple
