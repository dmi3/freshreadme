# freshReadme

Keep source code examples in your readme fresh!

## What is it?

This is simple Markdown preprocessor which updates examples from source code files.

If included files are compiled and verified on project build phase, it will help avoid a shameful situation when examples in readme are outdated.

Suggested usage is Git [pre commit hook](#automation).

## Motivation

Benefits comparing with [existing solutions](#alternative-solutions):

* Examples are updated in Markdown file not in resulting html, so default readme on Github or people who open Markdown file will get up to date examples
* Source is refreshed in one Markdown file, so you don't need to keep duplicated `README.source.md` and `README.source.md`
* Allows to include *only part* of souce file, to show only meaningful part of code skipping initialization and verification
* Includes source surrounded by specially formatted comments, not line numbers so if you change source code file, you don't need update includes in readme
* Source language agnostic

## How to use

* In `README.md`:
    - Put `<!-- [freshReadmeSource](path/to/file.ext) -->` before &grave;&grave;&grave; to inlude `path/to/file.ext`
    - Put `<!-- [freshReadmeSource](path/to/file.ext#snippetName) -->` before &grave;&grave;&grave; to inlude text surrounded by `snippetName` in `path/to/file.ext`

## Example

For example `README.md`:

<!-- [freshReadmeSource](examples/README.md.example) -->
```markdown
    # Regular expressions

    In Java

    <!-- [freshReadmeSource](examples/Examples.java#snippet1) -->
    ```java
    // This will be overwritten from examples/Examples.java
    ```

    In Python

    <!-- [freshReadmeSource](examples/examples.py#snippet2) -->
    ```python
    # This will be overwritten from examples/examples.py
    ```
```

And `examples/java.java`:

<!-- [freshReadmeSource](examples/Examples.java) -->
```java
import java.util.ArrayList;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import static org.assertj.core.api.Assertions.assertThat;

public class Examples {
    @org.junit.Test
    public void example1() throws Exception {

        // freshReadmeSnippet: snippet1
        Pattern pattern = Pattern.compile("[a-z]+");
        Matcher matcher = pattern.matcher("abc cde fgf");

        ArrayList<String> matches = new ArrayList<>();
        while (matcher.find()) {
            matches.add(matcher.group());
        }
        // freshReadmeSnippet: snippet1

        assertThat(matches).containsExactly("abc", "cde", "fgf");

    }
}

```


## Automation

Install freshReadme, then create `.git/hooks/pre-commit`:

<!-- [freshReadmeSource](examples/pre-commit) -->
```sh
#!/bin/sh

freshReadme
git add README.md
```

## Alternative solutions

* [okreadme](https://github.com/wan2land/okreadme) requires to keep duplicate files - `README.ok.md` and `README.md`
* [Apache Maven Site Plugin](https://maven.apache.org/guides/mini/guide-snippet-macro.html) works only in Java and mvn ecosystem, allows including files to resulting html
