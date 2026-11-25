package hr.foi.air.sevenpay.core.model

data class ImportResult(
    val imported: Int,
    val errors: List<RowError>?
)

data class RowError(
    val row: Int,
    val reason: String
)
