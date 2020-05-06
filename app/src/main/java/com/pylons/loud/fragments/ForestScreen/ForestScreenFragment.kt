package com.pylons.loud.fragments.ForestScreen

import android.os.Bundle
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup

import com.pylons.loud.R
import com.pylons.loud.fragments.PlayerAction.PlayerActionFragment
import com.pylons.loud.models.PlayerAction
import kotlinx.android.synthetic.main.fragment_forest_screen.*

// TODO: Rename parameter arguments, choose names that match
// the fragment initialization parameters, e.g. ARG_ITEM_NUMBER
private const val ARG_PARAM1 = "param1"
private const val ARG_PARAM2 = "param2"

/**
 * A simple [Fragment] subclass.
 * Use the [ForestScreenFragment.newInstance] factory method to
 * create an instance of this fragment.
 */
class ForestScreenFragment : Fragment() {
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
        return inflater.inflate(R.layout.fragment_forest_screen, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        text_forest_screen.setText(R.string.forest_desc)
        val frag = childFragmentManager.findFragmentById(R.id.fragment_player_action) as PlayerActionFragment
        frag.setAdapter(listOf(
            PlayerAction(1, getString(R.string.rabbit)),
            PlayerAction(2, getString(R.string.goblin)),
            PlayerAction(3, getString(R.string.wolf)),
            PlayerAction(4, getString(R.string.troll)),
            PlayerAction(5, getString(R.string.giant))
        ))
    }

    companion object {
        /**
         * Use this factory method to create a new instance of
         * this fragment using the provided parameters.
         *
         * @param param1 Parameter 1.
         * @param param2 Parameter 2.
         * @return A new instance of fragment ForestScreenFragment.
         */
        // TODO: Rename and change types and number of parameters
        @JvmStatic
        fun newInstance(param1: String, param2: String) =
            ForestScreenFragment().apply {
                arguments = Bundle().apply {
                    putString(ARG_PARAM1, param1)
                    putString(ARG_PARAM2, param2)
                }
            }
    }
}