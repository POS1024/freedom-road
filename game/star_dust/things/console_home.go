package things

import (
	"star_dust/consts"
	"star_dust/inter"
)

func HideHomeButtons(consoleButtons map[consts.ConsoleButtonKey]inter.Button) {
	inter.HideButtons(consoleButtons, []consts.ConsoleButtonKey{
		consts.ConsoleButtonNewGame,
		consts.ConsoleButtonContinueGame,
		consts.ConsoleButtonSetting,
		consts.ConsoleButtonQuit,
	})
}

func RefreshHomeButtons(consoleButtons map[consts.ConsoleButtonKey]inter.Button) {
	inter.ShowButtons(consoleButtons, []consts.ConsoleButtonKey{
		consts.ConsoleButtonNewGame,
		consts.ConsoleButtonContinueGame,
		consts.ConsoleButtonSetting,
		consts.ConsoleButtonQuit,
	})
}

func InitHomeButtons(consoleButtons map[consts.ConsoleButtonKey]inter.Button, manager *ConsoleManager) {
	consoleButtons[consts.ConsoleButtonNewGame] = newConsoleButton(50, 350, 200, 40, "New Game", true, func() error {
		HideHomeButtons(consoleButtons)
		manager.HideBackground()
		SM.Show()
		RemakeMainCharacters()
		ShowMCS()
		RemakeEnemyCharacters()
		ShowECS()
		RefreshGameButtons(consoleButtons)
		return nil
	})
	consoleButtons[consts.ConsoleButtonContinueGame] = newConsoleButton(50, 390, 200, 40, "Continue Game", true, func() error {
		HideHomeButtons(consoleButtons)
		manager.HideBackground()
		SM.Show()
		ShowMCS()
		ShowECS()
		RefreshGameButtons(consoleButtons)
		return nil
	})
	consoleButtons[consts.ConsoleButtonSetting] = newConsoleButton(50, 430, 200, 40, "Setting", true, func() error {
		HideHomeButtons(consoleButtons)
		RefreshSettingButtons(consoleButtons)
		return nil
	})
	consoleButtons[consts.ConsoleButtonQuit] = newConsoleButton(50, 470, 200, 40, "Quit", true, func() error {
		consts.ClosingCommand()
		return nil
	})
}
