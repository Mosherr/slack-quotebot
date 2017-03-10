// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package quotebot

var (
	// Set the token variable. It is needed to verify that the
	// requests to the slash command come from Slack. It is provided for
	// you by Slack when you create the Slash command as a custom
	// integration. https://my.slack.com/services/new/slash-commands
	token string = "QtFJq3lpxo491tdjPgJiaKRI"

	quotes = []string{
		"@gardak: Hans smells.",
		"@gardak: Greg likes men.",
		"@dafopp: I like it in the butt.",
	}
)