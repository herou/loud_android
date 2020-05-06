package com.pylons.loud.fragments.PylonCentralScreen

import android.os.Bundle
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup

import com.pylons.loud.R
import com.pylons.loud.fragments.PlayerAction.PlayerActionFragment
import com.pylons.loud.models.PlayerAction
import kotlinx.android.synthetic.main.fragment_pylon_central.*
import java.util.logging.Logger

// TODO: Rename parameter arguments, choose names that match
// the fragment initialization parameters, e.g. ARG_ITEM_NUMBER
private const val ARG_PARAM1 = "param1"
private const val ARG_PARAM2 = "param2"

/**
 * A simple [Fragment] subclass.
 * Use the [PylonCentralFragment.newInstance] factory method to
 * create an instance of this fragment.
 */
class PylonCentralFragment : Fragment() {
    private val Log = Logger.getLogger(PylonCentralFragment::class.java.name)

    // TODO: Rename and change types of parameters
    private var param1: String? = null
    private var param2: String? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        arguments?.let {
            param1 = it.getString(ARG_PARAM1)
            param2 = it.getString(ARG_PARAM2)
        }
    }

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        // Inflate the layout for this fragment
        return inflater.inflate(R.layout.fragment_pylon_central, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        text_pylon_central.setText(R.string.pylons_central_desc)

        val frag = childFragmentManager.findFragmentById(R.id.fragment_player_action) as PlayerActionFragment
        frag.setAdapter(listOf(
            PlayerAction(1, getString(R.string.buy_characters)),
            PlayerAction(2, getString(R.string.buy_5000_with_100_pylons)),
            PlayerAction(
                3,
                getString(R.string.sell_gold_from_orderbook_place_order_to_buy)
            ),
            PlayerAction(
                4,
                getString(R.string.buy_gold_from_orderbook_place_order_to_sell)
            ),
            PlayerAction(
                5,
                getString(R.string.sell_weapon_from_orderbook_place_order_to_buy)
            ),
            PlayerAction(
                6,
                getString(R.string.buy_weapon_from_orderbook_place_order_to_sell)
            ),
            PlayerAction(
                7,
                getString(R.string.sell_character_from_orderbook_place_order_to_buy)
            ),
            PlayerAction(
                8,
                getString(R.string.buy_character_from_orderbook_place_order_to_sell)
            ),
            PlayerAction(9, getString(R.string.update_character_name))
        ))
    }


    companion object {
        /**
         * Use this factory method to create a new instance of
         * this fragment using the provided parameters.
         *
         * @param param1 Parameter 1.
         * @param param2 Parameter 2.
         * @return A new instance of fragment PylonCentralFragment.
         */
        // TODO: Rename and change types and number of parameters
        @JvmStatic
        fun newInstance(param1: String, param2: String) =
            PylonCentralFragment().apply {
                arguments = Bundle().apply {
                    putString(ARG_PARAM1, param1)
                    putString(ARG_PARAM2, param2)
                }
            }
    }
}