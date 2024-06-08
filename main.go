/**
 * Fiber Admin
 *
 * Copyright [2024] [6GOD]
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import "fiber-admin/cmd"

//	@title			Fiber Admin Server API Documentation
//	@version		1.0
//	@description	This is the API documentation for the Fiber Admin Server.
//	@author			6GODD

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0

//	@contact.name	6GODD (BoChen SHEN)
//	@contact.email	6goddddddd@gmail.com | shenbochennnn@gmail.com
//	@host			localhost:3000
//	@BasePath		/api/v1

func main() {
	cmd.Execute()
}
