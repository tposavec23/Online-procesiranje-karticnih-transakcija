package hr.foi.air.sevenpay.ui.catalog

import hr.foi.air.sevenpay.core.model.Article
import java.math.BigDecimal

val sampleArticles = listOf(
    Article(1, "S-001", "Čokoladna torta", "Kriška, tamna čokolada i vrhnje", BigDecimal("4.50"), 8),
    Article(2, "S-002", "Sacher kocka", "Marelica i čokoladna glazura", BigDecimal("3.80"), 12),
    Article(3, "S-003", "Krempita", "Domaća, s vanilija kremom", BigDecimal("3.20"), 0),
    Article(4, "S-004", "Palačinka s Nutellom", "2 komada, posipana lješnjacima", BigDecimal("3.00"), 2)
)