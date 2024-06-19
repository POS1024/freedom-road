package inter

type DamageableItems interface {
	DamageCalculation(damage int, face int) (die bool, ex int)
}
