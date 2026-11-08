
package hr.foi.air.sevenpay.ui.theme

import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color

private val LightColorScheme = lightColorScheme(
    primary = Raspberry,
    onPrimary = Color.White,

    primaryContainer = LightRose,
    onPrimaryContainer = DarkRaspberry,

    secondary = DarkRaspberry,
    onSecondary = Color.White,

    tertiary = HeroPeach,
    onTertiary = PrimaryText,

    background = CreamBackground,
    onBackground = PrimaryText,

    surface = CreamBackground,
    onSurface = PrimaryText,

    surfaceVariant = LightRose,
    onSurfaceVariant = SecondaryText,

    outline = Border,
    error = Color(0xFFB3261E)
)

@Composable
fun SevenPayTheme(
    content: @Composable () -> Unit
) {
    MaterialTheme(
        colorScheme = LightColorScheme,
        typography = Typography,
        content = content
    )
}
