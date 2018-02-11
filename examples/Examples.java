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
