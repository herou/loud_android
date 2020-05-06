package com.pylons.loud.models

class Character(
    val id: String,
    val name: String,
    var level: Int,
    var price: Int,
    var xp: Double,
    var hp: Int,
    var maxHP: Int,
    var giantKill: Int,
    var lastUpdate: Int
) {

}