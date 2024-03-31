---
title: Advanced JSON customization in Go
posted_at: 25 03 2024
slug: Advanced_json_customization_in_go
tldr: Learn how you can use the Marshaler and Unmarshaler interfaces to customize encoding and decoding JSON
draft: false
---

Go interfaces are awesome, in this article I will show you how you can use the json.Marshaler and json.Unmarshaler interfaces to customize 
encoding and decoding JSON.
I assume you already have a basic knowledge of how to encode and decode in Go so I won't go over these. Let's get straight to the point.

## Encoding
Every data type in Go except **chan** **func** and **complex** has it's corresponding JSON type. But it can happen that you don't want to use the default
JSON type in which your struct, map or field will be encoded into. Fortunately, Go provides a way to customize the encoding of all supported data types.

By default when encoding a type, Go looks to see if that type implements the json.Marshaler interface
```go
type Marshaler interface {
    MarshalJSON()([]byte, error)
}
```
If the type implements that interface then Go will call the MarshalJSON function and use the returned byte slice as the encoded JSON value.
If the type does not implement that interface then Go will just fall back to trying to encode it based on it's own set of internal rules.

We can see this in action by looking at the `time.Time` type. Under the hood, it is just a struct 
```go
type Time struct {
	wall uint64
	ext  int64
	loc *Location
}
```
and it implements the Marshaler interface 
```go
func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(RFC3339Nano)+len(`""`))
	b = append(b, '"')
	b, err := t.appendStrictRFC3339(b)
	b = append(b, '"')
	if err != nil {
		return nil, errors.New("Time.MarshalJSON: " + err.Error())
	}
	return b, nil
}
```
So everytime Go encodes a time.Time field, it will call this function.

I will now demonstrate how you can do this using a very simple example and you can adapt it to meet your needs.

Let's say we are building an app that lets users send messages using carrier pigeons. They can browse a list of pigeons and pick one they like to be the deliverer of their next message.

Let's define a struct to hold our pigeon data
```go
type Pigeon struct {
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Experience int    `json:"experience"`
	Rating     int    `json:"rating"`
	Booked     bool   `json:"booked"`
}
```
Now lets assume we also want to display the pigeon's age in both human and pigeon years, but the secret formula to compute the pigeon's age in human years is only to be known by the server.
That means when encoding the pigeons data we have to change the way we encode the Age and Experience field to show something like this "pigeon_age years (human_years in human years)"
if that makes sense xD

Let's do that. First we need to create a custom type for the Age and Experience fields. I will call it PigeonYear and it will just be an alias for the int type
```go
type PigeonYear int
```
Then we have our secret formula to convert pigeon age to human years
```go
func pigeonAgeToHumanYears(pigeonYears int) int {
	if pigeonYears <= 1 {
		return pigeonYears * 5
	}
	return 5 + (pigeonYears-1)*2 
}
```
Now we have to make our PigeonYear type satisfy the Marshaler interface by implementing the MarshalJSON method
```go
func (p PigeonYear) MarshalJSON() ([]byte, error) {
	humanYears := pigeonAgeToHumanYears(int(p))
    //This outputs the raw string without quotes
	jsonValue := fmt.Sprintf("%d years old (%d in human years)", p, humanYears)
    //We use this to surround it with quotes so it can be a valid json string
	quotedJsonValue := strconv.Quote(jsonValue)
	return []byte(quotedJsonValue), nil
}
```
We also have to update our Pigeon struct to now use PigeonYear instead of int as type for the Age and Experience fields
```go
type Pigeon struct {
	Name       string     `json:"name"`
	Age        PigeonYear `json:"age"`
	Experience PigeonYear `json:"experience"`
	Rating     int        `json:"rating"`
	Booked     bool       `json:"booked"`
}
```
And now when we create a new instance of the pigeon struct and marshal it we get this as result
```go
func main() {
	Percy := Pigeon{
		Name:       "Percy",
		Age:        3,
		Experience: 2,
		Rating:     4,
		Booked:     true,
	}
	data, err := json.MarshalIndent(Percy, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
```
Notice that I used **MarshalIndent** instead of Marshal. This is just for readability reasons. You should not be using unless you need to print the JSON
in a human readable way because it makes the produces byte slice larger and also has a notable impact on performance because of the extra work Go has to
do to add whitespaces.

It gives us the following output
```sh
{
  "name": "Percy",
  "age": "3 years old (9 in human years)",
  "experience": "2 years old (7 in human years)",
  "rating": 4,
  "booked": true
}
```
So as you can see the rating field which is just an int pretty much as our age and experience fields, was encoded as an int, but because the PigeonYear type is a Marshaler,
Go used it's MarshalJSON field instead.

## Decoding
The same example goes for decoding. If the type that is being decoded implements the Unmarshaler interface
```go
type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}
```
Go calls it's UnmarshalJSON function when decoding it. Let's look at an example.

Now we need to add new pigeons to the database, so we send a POST request from our admin dashboard to our server and it needs to unmarshal the body into a pigeon.
We need to parse the int sent by the client into a PigeonYear, so lets implement the UnmarshalJSON method on our PigeonYear type
```go
func (p *PigeonYear) UnmarshalJSON(data []byte) error {
	var intData int
	if err := json.Unmarshal(data, &intData); err != nil {
		return err
	}
	*p = PigeonYear(intData)
	return nil
}
```
Note how we used a pointer receiver method, because in this case we need to be able to alter the value of the type.
We can test it with this sample code
```go
func main() {
	pigeonJSON := []byte(`{"name":"Percy","age":3,"experience":2,"rating":4,"booked":true}`)
	var pigeon Pigeon
	if err := json.Unmarshal(pigeonJSON, &pigeon); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	fmt.Println("Decoded Pigeon:", pigeon)
}
```
Which outputs
```sh
Decoded Pigeon: {Percy 3 2 4 true}
```

That's it. I hope you learned something new in this article. You can find all the code used [here](https://gist.github.com/TheWisePigeon/d40d1f5ec73dd9bca415600db3fba70b) 

Until next time :)
