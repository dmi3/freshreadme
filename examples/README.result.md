# Regular expressions

In Java

<!-- [freshReadmeSource](examples/Examples.java#snippet1) -->
```java
        Pattern pattern = Pattern.compile("[a-z]+");
        Matcher matcher = pattern.matcher("abc cde fgf");

        ArrayList<String> matches = new ArrayList<>();
        while (matcher.find()) {
            matches.add(matcher.group());
        }
```

In Python

<!-- [freshReadmeSource](examples/examples.py#snippet2) -->
```python
matches = re.findall("[a-z]+", "abc cde fgf")
```
