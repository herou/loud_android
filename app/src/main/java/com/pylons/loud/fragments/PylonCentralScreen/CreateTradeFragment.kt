package com.pylons.loud.fragments.PylonCentralScreen

import android.content.Context
import android.os.Bundle
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Toast
import androidx.appcompat.app.AlertDialog
import androidx.fragment.app.activityViewModels
import androidx.lifecycle.Observer
import androidx.recyclerview.widget.RecyclerView

import com.pylons.loud.R
import com.pylons.loud.activities.GameScreenActivity
import com.pylons.loud.constants.Coin
import com.pylons.loud.constants.Item.ANGEL_SWORD
import com.pylons.loud.constants.Item.BRONZE_SWORD
import com.pylons.loud.constants.Item.COPPER_SWORD
import com.pylons.loud.constants.Item.DROP_DRAGONACID
import com.pylons.loud.constants.Item.DROP_DRAGONFIRE
import com.pylons.loud.constants.Item.DROP_DRAGONICE
import com.pylons.loud.constants.Item.GOBLIN_EAR
import com.pylons.loud.constants.Item.IRON_SWORD
import com.pylons.loud.constants.Item.SILVER_SWORD
import com.pylons.loud.constants.Item.TROLL_TOES
import com.pylons.loud.constants.Item.WOLF_TAIL
import com.pylons.loud.constants.Item.WOODEN_SWORD
import com.pylons.loud.constants.Trade
import com.pylons.loud.fragments.Character.CharacterFragment
import com.pylons.loud.fragments.Character.MyCharacterRecyclerViewAdapter
import com.pylons.loud.fragments.Item.ItemFragment
import com.pylons.loud.fragments.Item.MyItemRecyclerViewAdapter
import com.pylons.loud.fragments.itemspec.ItemSpecFragment
import com.pylons.loud.fragments.itemspec.MyItemSpecRecyclerViewAdapter
import com.pylons.loud.models.Character
import com.pylons.loud.models.Item
import com.pylons.loud.models.trade.CharacterSpec
import com.pylons.loud.models.trade.MaterialSpec
import com.pylons.loud.models.trade.Spec
import com.pylons.loud.models.trade.WeaponSpec
import com.pylons.wallet.core.types.tx.recipe.*
import kotlinx.android.synthetic.main.create_trade_buy.*
import kotlinx.android.synthetic.main.create_trade_confirm.*
import kotlinx.android.synthetic.main.create_trade_sell.*
import kotlinx.android.synthetic.main.dialog_input.view.*
import kotlinx.android.synthetic.main.fragment_create_trade.*
import java.util.logging.Logger

/**
 * A simple [Fragment] subclass.
 */
class CreateTradeFragment : Fragment() {
    private val Log = Logger.getLogger(CreateTradeFragment::class.java.name)

    private var listener: OnFragmentInteractionListener? = null

    private var coinInput = listOf<CoinInput>()
    private var itemInput = listOf<ItemInput>()
    private var coinOutput = listOf<CoinOutput>()
    private var itemOutput = listOf<com.pylons.wallet.core.types.tx.item.Item>()
    private var extraInfo = Trade.DEFAULT
    private lateinit var itemBuyFragment: ItemSpecFragment
    private lateinit var characterSellFragment: CharacterFragment
    private lateinit var itemSellFragment: ItemFragment

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        // Inflate the layout for this fragment
        return inflater.inflate(R.layout.fragment_create_trade, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        initFragments()
        initModel()
        initTextPylonsBuy()
        initTextGoldBuy()
        initCharacterBuy()
        initItemBuy()
        initTextGoldSell()
        initCharacterSell()
        initItemSell()
        initCreateTradeButton()
    }

    override fun onAttach(context: Context) {
        super.onAttach(context)
        if (context is OnFragmentInteractionListener) {
            listener = context
        } else {
            throw RuntimeException(context.toString() + " must implement OnListFragmentInteractionListener")
        }
    }

    override fun onDetach() {
        super.onDetach()
        listener = null
    }

    interface OnFragmentInteractionListener {
        fun onCreateTrade(
            coinInput: List<CoinInput>,
            itemInput: List<ItemInput>,
            coinOutput: List<CoinOutput>,
            itemOutput: List<com.pylons.wallet.core.types.tx.item.Item>,
            extraInfo: String
        )
    }

    private fun initFragments() {
        itemBuyFragment =
            childFragmentManager.findFragmentById(R.id.fragment_item_buy) as ItemSpecFragment
        characterSellFragment =
            childFragmentManager.findFragmentById(R.id.fragment_character_sell) as CharacterFragment
        itemSellFragment =
            childFragmentManager.findFragmentById(R.id.fragment_item_sell) as ItemFragment
        childFragmentManager.beginTransaction().hide(itemBuyFragment).commit()
        childFragmentManager.beginTransaction().hide(characterSellFragment).commit()
        childFragmentManager.beginTransaction().hide(itemSellFragment).commit()
    }

    private fun initModel() {
        val model: GameScreenActivity.SharedViewModel by activityViewModels()
        model.getTradeInput().observe(viewLifecycleOwner, Observer { item ->
            if (item != null) {
                Log.info(item.toString())
                itemInput = when (item) {
                    is CharacterSpec ->
                        listOf(
                            ItemInput(
                                listOf(
                                    DoubleInputParam(
                                        "XP",
                                        item.xp.min.toString(),
                                        item.xp.max.toString()
                                    )
                                ),
                                listOf(
                                    LongInputParam(
                                        "level",
                                        item.level.min.toLong(),
                                        item.level.max.toLong()
                                    )
                                ),
                                listOf(StringInputParam("Name", item.name))
                            )
                        )

                    is WeaponSpec ->
                        listOf(
                            ItemInput(
                                listOf(
                                    DoubleInputParam(
                                        "attack",
                                        item.attack.min.toString(),
                                        item.attack.max.toString()
                                    )
                                ),
                                listOf(
                                    LongInputParam(
                                        "level",
                                        item.level.min.toLong(),
                                        item.level.max.toLong()
                                    )
                                ),
                                listOf(StringInputParam("Name", item.name))
                            )
                        )

                    is MaterialSpec ->
                        listOf(
                            ItemInput(
                                listOf(),
                                listOf(
                                    LongInputParam(
                                        "level",
                                        item.level.min.toLong(),
                                        item.level.max.toLong()
                                    )
                                ),
                                listOf(StringInputParam("Name", item.name))
                            )
                        )
                    else -> listOf()
                }
                extraInfo = when (item) {
                    is CharacterSpec -> Trade.CHARACTER_BUY
                    else -> Trade.ITEM_BUY
                }

                promptBuyOrderStep2()
                model.setTradeInput(null)
            }
        })

        model.getTradeOutput().observe(viewLifecycleOwner, Observer { output ->
            if (output != null) {
                Log.info(output.toString())

                val c = context
                if (c != null) {
                    itemOutput = listOf(output)
                    extraInfo = if (output.strings["Type"] == "Character") {
                        Trade.CHARACTER_SELL
                    } else {
                        Trade.ITEM_SELL
                    }
                    Log.info(itemOutput.toString())
                    displayConfirmTrade()
                    model.setTradeOutput(null)
                }
            }
        })
    }

    private fun initTextPylonsBuy() {
        text_pylons_buy.setOnClickListener {
            childFragmentManager.beginTransaction().hide(itemBuyFragment).commit()

            val c = context
            if (c != null) {
                val mDialogView = LayoutInflater.from(c).inflate(R.layout.dialog_input, null)

                val dialogBuilder = AlertDialog.Builder(c, R.style.MyDialogTheme)
                dialogBuilder.setMessage(
                    getString(R.string.trade_pylon_buy)
                )
                    .setCancelable(false)
                    .setPositiveButton("Proceed") { _, _ ->
                    }
                    .setNegativeButton("Cancel") { dialog, _ ->
                        dialog.cancel()
                    }

                val alert = dialogBuilder.create()
                alert.setTitle("Confirm")
                alert.setView(mDialogView)
                alert.show()

                alert.getButton(AlertDialog.BUTTON_POSITIVE).setOnClickListener {
                    if (mDialogView.edit_text_amount.text.toString() != "" && mDialogView.edit_text_amount.text.toString()
                            .toLong() > 0
                    ) {
                        coinInput = listOf(
                            CoinInput(
                                Coin.PYLON,
                                mDialogView.edit_text_amount.text.toString().toLong()
                            )
                        )

                        displaySellTrade()
                        alert.dismiss()
                    } else {
                        Toast.makeText(
                            c,
                            getString(R.string.enter_valid_amount),
                            Toast.LENGTH_SHORT
                        ).show()
                    }
                }
            }
        }
    }

    private fun initTextGoldBuy() {
        text_gold_buy.setOnClickListener {
            childFragmentManager.beginTransaction().hide(itemBuyFragment).commit()

            val c = context
            if (c != null) {
                val mDialogView = LayoutInflater.from(c).inflate(R.layout.dialog_input, null)

                val dialogBuilder = AlertDialog.Builder(c, R.style.MyDialogTheme)
                dialogBuilder.setMessage(
                    getString(R.string.trade_gold_buy)
                )
                    .setCancelable(false)
                    .setPositiveButton("Proceed") { _, _ ->
                    }
                    .setNegativeButton("Cancel") { dialog, _ ->
                        dialog.cancel()
                    }

                val alert = dialogBuilder.create()
                alert.setTitle("Confirm")
                alert.setView(mDialogView)
                alert.show()

                alert.getButton(AlertDialog.BUTTON_POSITIVE).setOnClickListener {
                    if (mDialogView.edit_text_amount.text.toString() != "" && mDialogView.edit_text_amount.text.toString()
                            .toLong() > 0
                    ) {
                        coinInput = listOf(
                            CoinInput(
                                Coin.LOUD,
                                mDialogView.edit_text_amount.text.toString().toLong()
                            )
                        )
                        extraInfo = Trade.DEFAULT
                        promptBuyOrderStep2()
                        alert.dismiss()
                    } else {
                        Toast.makeText(
                            c,
                            getString(R.string.enter_valid_amount),
                            Toast.LENGTH_SHORT
                        ).show()
                    }
                }
            }
        }
    }

    private fun initCharacterBuy() {
        text_character_buy.setOnClickListener {
            val adapter = MyItemSpecRecyclerViewAdapter(
                listOf(
                    CharacterSpec("Lion", Spec(1, 2), Spec(1.0, 1000.0)),
                    CharacterSpec("Liger", Spec(2, 1000), Spec(1.0, 1000.0))
                ),
                itemBuyFragment.getListener()
            )
            val view1 = itemBuyFragment.view as RecyclerView
            view1.adapter = adapter
            childFragmentManager.beginTransaction().show(itemBuyFragment).commit()
        }
    }

    private fun initItemBuy() {
        text_item_buy.setOnClickListener {
            val adapter = MyItemSpecRecyclerViewAdapter(
                listOf(
                    WeaponSpec(
                        WOODEN_SWORD,
                        Spec(1, 1),
                        Spec(3, 3)
                    ),
                    WeaponSpec(
                        WOODEN_SWORD,
                        Spec(2, 2),
                        Spec(6, 6)
                    ),
                    WeaponSpec(
                        COPPER_SWORD,
                        Spec(1, 1),
                        Spec(10, 10)
                    ),
                    WeaponSpec(
                        COPPER_SWORD,
                        Spec(2, 2),
                        Spec(20, 20)
                    ),
                    WeaponSpec(
                        SILVER_SWORD,
                        Spec(1, 1),
                        Spec(30, 30)
                    ),
                    WeaponSpec(
                        BRONZE_SWORD,
                        Spec(1, 1),
                        Spec(50, 50)
                    ),
                    WeaponSpec(
                        IRON_SWORD,
                        Spec(1, 1),
                        Spec(100, 100)
                    ),
                    WeaponSpec(
                        ANGEL_SWORD,
                        Spec(1, 1),
                        Spec(1000, 1000)
                    ),
                    MaterialSpec(
                        TROLL_TOES,
                        Spec(1, 1)
                    ),
                    MaterialSpec(
                        WOLF_TAIL,
                        Spec(1, 1)
                    ),
                    MaterialSpec(
                        GOBLIN_EAR,
                        Spec(1, 1)
                    ),
                    MaterialSpec(
                        TROLL_TOES,
                        Spec(1, 1)
                    ),
                    MaterialSpec(
                        DROP_DRAGONFIRE,
                        Spec(1, 1)
                    ),
                    MaterialSpec(
                        DROP_DRAGONICE,
                        Spec(1, 1)
                    ),
                    MaterialSpec(
                        DROP_DRAGONACID,
                        Spec(1, 1)
                    )
                ), itemBuyFragment.getListener()
            )
            val myView = itemBuyFragment.view as RecyclerView
            myView.adapter = adapter
            childFragmentManager.beginTransaction().show(itemBuyFragment).commit()
        }
    }

    private fun promptBuyOrderStep2() {
        val c = context
        if (c != null) {
            val mDialogView = LayoutInflater.from(c).inflate(R.layout.dialog_input, null)
            val dialogBuilder = AlertDialog.Builder(c, R.style.MyDialogTheme)
            dialogBuilder.setMessage(
                getString(R.string.trade_pylon_offer)
            )
                .setCancelable(false)
                .setPositiveButton("Proceed") { _, _ ->
                }
                .setNegativeButton("Cancel") { dialog, _ ->
                    dialog.cancel()
                }

            val alert = dialogBuilder.create()
            alert.setTitle("Confirm")
            alert.setView(mDialogView)
            alert.show()

            alert.getButton(AlertDialog.BUTTON_POSITIVE).setOnClickListener {
                if (mDialogView.edit_text_amount.text.toString() != "" && mDialogView.edit_text_amount.text.toString()
                        .toLong() > 0
                ) {
                    coinOutput = listOf(
                        CoinOutput(
                            Coin.PYLON,
                            mDialogView.edit_text_amount.text.toString().toLong()
                        )
                    )
                    displayConfirmTrade()
                    alert.dismiss()
                } else {
                    Toast.makeText(c, getString(R.string.enter_valid_amount), Toast.LENGTH_SHORT)
                        .show()
                }
            }
        }
    }

    private fun displaySellTrade() {
        layout_create_trade_buy.visibility = View.GONE
        layout_create_trade_sell.visibility = View.VISIBLE
        layout_create_trade_confirm.visibility = View.GONE
    }

    private fun initTextGoldSell() {
        text_gold_sell.setOnClickListener {
            with(childFragmentManager) {
                beginTransaction().hide(characterSellFragment).commit()
                beginTransaction().hide(itemSellFragment).commit()
            }

            val c = context
            if (c != null) {
                val mDialogView = LayoutInflater.from(c).inflate(R.layout.dialog_input, null)

                val dialogBuilder = AlertDialog.Builder(c, R.style.MyDialogTheme)
                dialogBuilder.setMessage(
                    getString(R.string.trade_gold_sell)
                )
                    .setCancelable(false)
                    .setPositiveButton("Proceed") { _, _ ->
                    }
                    .setNegativeButton("Cancel") { dialog, _ ->
                        dialog.cancel()
                    }

                val alert = dialogBuilder.create()
                alert.setTitle("Confirm")
                alert.setView(mDialogView)
                alert.show()

                alert.getButton(AlertDialog.BUTTON_POSITIVE).setOnClickListener {
                    if (mDialogView.edit_text_amount.text.toString() != "" && mDialogView.edit_text_amount.text.toString()
                            .toLong() > 0
                    ) {
                        coinOutput = listOf(
                            CoinOutput(
                                Coin.LOUD,
                                mDialogView.edit_text_amount.text.toString().toLong()
                            )
                        )
                        extraInfo = Trade.DEFAULT
                        displayConfirmTrade()
                        alert.dismiss()
                    } else {
                        Toast.makeText(
                            c,
                            getString(R.string.enter_valid_amount),
                            Toast.LENGTH_SHORT
                        ).show()
                    }
                }
            }
        }
    }

    private fun initCharacterSell() {
        val model: GameScreenActivity.SharedViewModel by activityViewModels()
        val player = model.getPlayer().value
        val list = player?.characters ?: listOf<Character>()
        text_character_sell.text = "${getString(R.string.trade_character)} (${list.size})"

        text_character_sell.setOnClickListener {
            val player = model.getPlayer().value
            val list = player?.characters ?: listOf<Character>()

            if (list.isNotEmpty()) {
                val adapter = MyCharacterRecyclerViewAdapter(
                    list,
                    characterSellFragment.getListener(),
                    3
                )
                val view1 = characterSellFragment.view as RecyclerView
                view1.adapter = adapter

                with(childFragmentManager) {
                    beginTransaction().hide(itemSellFragment).commit()
                    beginTransaction().show(characterSellFragment).commit()
                }
            }
        }
    }

    private fun initItemSell() {
        val model: GameScreenActivity.SharedViewModel by activityViewModels()
        val player = model.getPlayer().value
        if (player != null) {
            val items = mutableListOf<Item>()
            items.addAll(player.weapons)
            items.addAll(player.materials)
            text_item_sell.text = "${getString(R.string.trade_item)} (${items.size})"
        }

        text_item_sell.setOnClickListener {
            val player = model.getPlayer().value
            if (player != null) {
                val items = mutableListOf<Item>()
                items.addAll(player.weapons)
                items.addAll(player.materials)

                if (items.isNotEmpty()) {
                    val adapter = MyItemRecyclerViewAdapter(items, itemSellFragment.getListener(), 5)
                    val myView = itemSellFragment.view as RecyclerView
                    myView.adapter = adapter

                    with(childFragmentManager) {
                        beginTransaction().hide(characterSellFragment).commit()
                        beginTransaction().show(itemSellFragment).commit()
                    }
                }
            }
        }
    }

    private fun displayConfirmTrade() {
        layout_create_trade_buy.visibility = View.GONE
        layout_create_trade_sell.visibility = View.GONE
        layout_create_trade_confirm.visibility = View.VISIBLE

        if (coinOutput.isNotEmpty()) {
            text_trade_output.text = "${coinOutput[0].amount} ${coinOutput[0].denom}"
        }

        if (coinInput.isNotEmpty()) {
            text_trade_input.text = "${coinInput[0].count} ${coinInput[0].coin}"
        }

        if (itemOutput.isNotEmpty()) {
            text_trade_output.text =
                "${itemOutput[0].strings["Name"]} Lv${itemOutput[0].longs["level"]}"
        }

        if (itemInput.isNotEmpty()) {
            text_trade_input.text =
                "${itemInput[0].strings.find { it.key == "Name" }?.value} Lv${itemInput[0].longs.find { it.key == "level" }?.maxValue}"
        }

    }

    private fun initCreateTradeButton() {
        button_create_trade.setOnClickListener {
            listener?.onCreateTrade(
                coinInput,
                itemInput,
                coinOutput,
                itemOutput,
                extraInfo
            )
        }
    }

}
