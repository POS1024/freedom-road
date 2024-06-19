package things

import (
	"star_dust/consts"
	"star_dust/inter"
)

func HideGameButtons(consoleButtons map[consts.ConsoleButtonKey]inter.Button) {
	inter.HideButtons(consoleButtons, []consts.ConsoleButtonKey{
		consts.ConsoleButtonGameBackToHome,
	})
}

func RefreshGameButtons(consoleButtons map[consts.ConsoleButtonKey]inter.Button) {
	inter.ShowButtons(consoleButtons, []consts.ConsoleButtonKey{
		consts.ConsoleButtonGameBackToHome,
	})
}

func InitGameButtons(consoleButtons map[consts.ConsoleButtonKey]inter.Button, manager *ConsoleManager) {
	consoleButtons[consts.ConsoleButtonGameBackToHome] = newConsoleButton(860, 60, 100, 40, "Exit", false, func() error {
		HideGameButtons(consoleButtons)
		SM.Hide()
		HideMCS()
		HideECS()
		manager.ShowBackground()
		RefreshHomeButtons(consoleButtons)
		return nil
	})
}
