import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.extension.ExtendWith;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.EnumSource;
import org.mockito.Mock;
import org.mockito.Mockito;
import org.mockito.junit.jupiter.MockitoExtension;

import java.time.DayOfWeek;
import static org.junit.jupiter.api.Assertions.*;

@ExtendWith(MockitoExtension.class)
public class BankOfficeTest {
    @Mock
    BankOffice bankOffice;

    @BeforeEach
    public void setupMock() {
        Mockito.when(bankOffice.daysOff()).thenReturn("SUNDAY");
    }

    @ParameterizedTest
    @EnumSource(value = DayOfWeek.class, names = {"SATURDAY", "SUNDAY"})
    @DisplayName("Day off is Sunday, not Saturday")
    public void testDaysOff(DayOfWeek day) {
        assertEquals(day.toString(), bankOffice.daysOff());
    }
}
