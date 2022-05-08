import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.assertEquals;

public class CodeTest {

    Code code = new Code();

    @Test
    public void testSayHello() {
        assertEquals("Hello world!", code.sayHello());
    }
}
