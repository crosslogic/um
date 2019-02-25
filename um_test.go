package um

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertir(t *testing.T) {
	var out float64
	var err error

	// Validacion
	assert.NotNil(t, UM{Tipo: "Ningun tipo", ID: "kg"}.Validar())
	assert.NotNil(t, UM{Tipo: TipoCantidad, ID: "sfds"}.Validar())
	assert.NotNil(t, UM{Tipo: TipoCantidad, ID: "sfds"}.Validar())
	assert.NotNil(t, UM{Tipo: TipoCantidad, ID: "kg"}.Validar())
	assert.Nil(t, UM{Tipo: TipoPeso, ID: "kg"}.Validar())

	// Conversiones entre distintos tipos
	out, err = Convertir(1, "t", "u")
	assert.NotNil(t, err)

	// Medidas inexistentes
	out, err = Convertir(1, "t", "xas")
	assert.NotNil(t, err)

	out, err = Convertir(1, "ter", "t")
	assert.NotNil(t, err)

	// Conversiones
	out, err = Convertir(1, "t", "g")
	assert.Equal(t, 1000000., out)
	assert.Nil(t, err)

	out, err = Convertir(1000, "g", "kg")
	assert.Equal(t, 1.0, out)
	assert.Nil(t, err)

	out, err = Convertir(1, "kg", "t")
	assert.Equal(t, 0.001, out)
	assert.Nil(t, err)

	out, err = Convertir(1, "u", "u")
	assert.Equal(t, 1., out)
	assert.Nil(t, err)

	out, err = Convertir(1, "m3", "L")
	assert.Equal(t, 1000., out)
	assert.Nil(t, err)

	out, err = Convertir(1, "ml", "cm3")
	assert.Equal(t, 1., out)
	assert.Nil(t, err)

	// Relaciones ingresadas
	out, err = Convertir(600, "ml", "kg", RelacionUM{
		Un:        "L",
		EquivaleA: 0.92,
		De:        "kg"},
	)
	assert.Nil(t, err)
	assert.Equal(t, 0.552, out)

	// Inversa
	out, err = Convertir(600, "ml", "kg", RelacionUM{
		Un:        "kg",
		EquivaleA: 0.92,
		De:        "L"},
	)
	assert.Nil(t, err)
	assert.Equal(t, 0.6521739130434782, out)

	// Inversa
	out, err = Convertir(1, "kg", "L", RelacionUM{
		Un:        "kg",
		EquivaleA: 0.92,
		De:        "L"},
	)
	assert.Nil(t, err)
	assert.Equal(t, 0.92, out)

	// Relaciones ingresadas
	out, err = Convertir(600, "ml", "kg", RelacionUM{
		Un:        "kg",
		EquivaleA: 0.92,
		De:        "cm"},
	)
	assert.NotNil(t, err)

	out, err = Convertir(2, "L", "g", RelacionUM{
		Un:        "L",
		EquivaleA: 1,
		De:        "kg"},
	)
	assert.Nil(t, err)
	assert.InDelta(t, 2000., out, 0.0000000001)

	// Desde y hasta no están en la relación
	out, err = Convertir(2, "L", "kg", RelacionUM{
		Un:        "ml",
		EquivaleA: 1,
		De:        "g"},
	)
	assert.Nil(t, err)
	assert.InDelta(t, 2., out, 0.0000000001)
}

func TestMedidasDe(t *testing.T) {

	{
		mm, err := MedidasDe("Cantidad")
		assert.Nil(t, err)
		assert.Len(t, mm, 1)
	}
	{
		mm, err := MedidasDe("Peso")
		assert.Nil(t, err)
		assert.Len(t, mm, 4)
	}

}
