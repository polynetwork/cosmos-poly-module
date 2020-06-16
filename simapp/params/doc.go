/*
 * Copyright (C) 2020 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

/*
Package params defines the simulation parameters in the simapp.

It contains the default weights used for each transaction used on the module's
simulation. These weights define the chance for a transaction to be simulated at
any gived operation.

You can repace the default values for the weights by providing a params.json
file with the weights defined for each of the transaction operations:

	{
		"op_weight_msg_send": 60,
		"op_weight_msg_delegate": 100,
	}

In the example above, the `MsgSend` has 60% chance to be simulated, while the
`MsgDelegate` will always be simulated.
*/
package params
