package things

import (
	"star_dust/config"
	"star_dust/consts"
	"star_dust/inter"
)

func HideSettingButtons(consoleButtons map[consts.ConsoleButtonKey]inter.Button) {
	inter.HideButtons(consoleButtons, []consts.ConsoleButtonKey{
		consts.ConsoleButtonSettingBackToHome,
		consts.ConsoleButtonSettingVolume,
		consts.ConsoleButtonSettingLanguage,
		consts.ConsoleButtonSettingSave,
	})
}

func RefreshSettingButtons(consoleButtons map[consts.ConsoleButtonKey]inter.Button) {
	inter.ShowButtons(consoleButtons, []consts.ConsoleButtonKey{
		consts.ConsoleButtonSettingBackToHome,
		consts.ConsoleButtonSettingVolume,
		consts.ConsoleButtonSettingLanguage,
		consts.ConsoleButtonSettingSave,
	})
}

func InitSettingButtons(consoleButtons map[consts.ConsoleButtonKey]inter.Button, manager *ConsoleManager) {
	consoleButtons[consts.ConsoleButtonSettingBackToHome] = newConsoleButton(50, 60, 100, 40, "Back", false, func() error {
		config.GameConfigurationConfigurator.Read()
		HideSettingButtons(consoleButtons)
		RefreshHomeButtons(consoleButtons)
		return nil
	})
	consoleButtons[consts.ConsoleButtonSettingSave] = newConsoleButton(860, 460, 100, 40, "Save", false, func() error {
		config.GameConfigurationConfigurator.Store()
		HideSettingButtons(consoleButtons)
		RefreshHomeButtons(consoleButtons)
		return nil
	})

	consoleButtons[consts.ConsoleButtonSettingVolume] = newConsoleButton(280, 160, 114, 40, "Volume", false, func() error {
		return nil
	})
	consoleButtons[consts.ConsoleButtonSettingVolume] = NewVolumeButton(280, 160, "Volume", false)

}
