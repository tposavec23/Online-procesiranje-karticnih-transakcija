
package hr.foi.air.sevenpay.ui.theme

import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color

private val SlatkoColorScheme = lightColorScheme(
    primary = Malina,
    onPrimary = Color.White,

    primaryContainer = SvijetlaRuza,
    onPrimaryContainer = TamnaMalina,

    secondary = TamnaMalina,
    onSecondary = Color.White,

    tertiary = HeroBreskva,
    onTertiary = PrimarniTekst,

    background = KremPozadina,
    onBackground = PrimarniTekst,

    surface = KremPozadina,
    onSurface = PrimarniTekst,

    surfaceVariant = SvijetlaRuza,
    onSurfaceVariant = SporedniTekst,

    outline = Rub,
    error = Color(0xFFB3261E)
)

@Composable
fun SevenPayTheme(
    content: @Composable () -> Unit
) {
    MaterialTheme(
        colorScheme = SlatkoColorScheme,
        typography = Typography,
        content = content
    )
}
