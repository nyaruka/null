v3.1.0 (2026-08-31)
-------------------------
 * Don't share backing arrays between JSON values scanned from NULL
 * Replace map contents when scanning instead of merging into them
 * Error rather than silently wrap when an integer doesn't fit the target type
 * Require Go 1.25 and drop golang.org/x/exp dependency

v3.0.0 (2023-09-06)
-------------------------
 * Convert null.Map to be generic

v2.0.3 (2023-05-10)
-------------------------
 * Clone the bytes driver gives us when scanning JSON

v2.0.2 (2023-02-01)
-------------------------
 * Add JSON.IsNull method

v2.0.1 (2023-02-01)
-------------------------
 * Export NullJSON constant

v2.0.0 (2023-01-31)
-------------------------
 * Int becomes Int64, Add Int for regular int
 * Generics to easily support custom int and string types
 * Simplify map so it's no longer a struct

