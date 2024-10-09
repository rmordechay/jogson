package tests

const jsonObjectTest = `{"name": "Jason", "Age": 15, "address": null, "height": 1.81, "is_funny": true}`
const jsonAnyArrayTest = `["Jason", 15, null, 1.81, true]`
const jsonStringArrayTest = `["Jason", "Chris", "Rachel"]`
const jsonIntArrayTest = `[0, 15, -54, -346, 9223372036854775807]`
const jsonFloatArrayTest = `[15.13, 2, 45.3984, -1.81, 9.223372036854776]`
const json2DIntArrayTest = `[[1, 2], [3, 4]]`
const json2DArrayTest = `[[1, 2], [3, 4], 3.23, null]`
const jsonObjectArrayTest = `[{"name": "Jason"}, {"name":  "Chris"}]`
const jsonObjectWithArrayTest = `{"names": ["Jason", "Chris", "Rachel"], "name": "Charlie", "address": null}`
const jsonObjectNestedArrayTest = `{"personTest": {"name": "Jason"}}`
const jsonArrayWithNullTest = `[{"name": "Jason"}, {"name": "Chris"}, "string", null]`
const jsonTimeTest = `{"time1": "2024-10-06T17:59:44Z", "time2": "2024-10-06T17:59:44+00:00", "time3": "Sunday, 06-Oct-24 17:59:44 UTC"}`
const jsonInvalidTimeTest = `{"time1": null, "time2": 0, "time3": "INVALID", "time4": false}`
