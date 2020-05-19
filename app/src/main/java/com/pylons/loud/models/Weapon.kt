package com.pylons.loud.models

import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class Weapon(
    override val id: String,
    override val name: String,
    override val level: Long,
    override val attack: Double,
    val price: Int,
    val preItem: List<String>,
    override val lastUpdate: Long
) : Item()
