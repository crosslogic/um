// Package um es para trabajar con unidades de medida.
package um

import (
	"github.com/pkg/errors"
)

var unidades map[string]UM

const (
	CubicMillimeter         = CubicMeter * 1e-9
	CubicCentimeter         = CubicMeter * 1e-6
	CubicDecimeter          = CubicMeter * 1e-3
	CubicMeter      float64 = 1e0
	CubicDecameter          = CubicMeter * 1e3
	CubicHectometer         = CubicMeter * 1e6
	CubicKilometer          = CubicMeter * 1e9
	Milliliter              = Liter * 1e-3
	Centiliter              = Liter * 1e-2
	Deciliter               = Liter * 1e-1
	Liter                   = CubicMeter * 1e-3
	Decaliter               = Liter * 1e1
	Hectoliter              = Liter * 1e2
	Kiloliter               = Liter * 1e3
)

const (
	Miligramo          = Gramo * 1e-3
	Centigramo         = Gramo * 1e-2
	Decigramo          = Gramo * 1e-1
	Gramo              = Kilogramo * 1e-3
	Decagramo          = Gramo * 1e1
	Hectogramo         = Gramo * 1e2
	Kilogramo  float64 = 1e0
	Tonelada           = Kilogramo * 1e3
)
const (
	Metro      float64 = 1e0
	Centimetro         = Metro * 1e-2
)

func init() {
	unidades = map[string]UM{
		// Cantidad
		"u": UM{"u", "Unidad", "Unidades", TipoCantidad, 1},

		// Peso
		"t":  UM{"t", "Tonelada", "Toneladas", TipoPeso, Tonelada},
		"kg": UM{"kg", "Kilogramo", "Kilogramos", TipoPeso, Kilogramo},
		"g":  UM{"g", "Gramo", "Gramos", TipoPeso, Gramo},
		"mg": UM{"mg", "Miligramo", "Miligramos", TipoPeso, Miligramo},

		// Volumen
		"m":  UM{"m", "Metro", "Metros", TipoDistancia, Metro},
		"cm": UM{"cm", "Centímetro", "Centímetro", TipoDistancia, Centimetro},

		// Volumen
		"km3":  UM{"km3", "Kilómetro cúbico", "Kilómetros cúbicos", TipoVolumen, CubicKilometer},
		"hm3":  UM{"hm3", "Hectómetro cúbico", "Hectómetros cúbicos", TipoVolumen, CubicHectometer},
		"dam3": UM{"dam3", "Decámetro cúbico", "Decámetros cúbicos", TipoVolumen, CubicDecameter},
		"m3":   UM{"m3", "Metro cúbico", "Metros cúbicos", TipoVolumen, CubicMeter},
		"dm3":  UM{"dm3", "Decímetro cúbico", "Decímetros cúbicos", TipoVolumen, CubicDecimeter},
		"cm3":  UM{"cm3", "Centímetro cúbico", "Centímetros cúbicos", TipoVolumen, CubicCentimeter},
		"mm3":  UM{"mm3", "Milímetro cúbico", "kilómetros cúbicos", TipoVolumen, CubicKilometer},

		// Capacidad
		"Kl": UM{"Kl", "Kilolitro", "Kilolitro", TipoVolumen, Kiloliter},
		"Hl": UM{"Hl", "Hectolitro", "Hectolitro", TipoVolumen, Hectoliter},
		"Dl": UM{"Dl", "Decalitro", "Decalitros", TipoVolumen, Decaliter},
		"L":  UM{"L", "Litro", "Litros", TipoVolumen, Liter},
		"dl": UM{"dl", "Decilitro", "Decilitros", TipoVolumen, Deciliter},
		"cl": UM{"cl", "Centilitro", "Centilitros", TipoVolumen, Centiliter},
		"ml": UM{"ml", "Mililitro", "Mililitros", TipoVolumen, Milliliter},
	}
}

// UMVolumen es una definición de una unidad de medida
type UM struct {
	ID           string
	Nombre       string
	NombrePlural string
	Tipo         string
	factor       float64
}

// NewUM devuelve una unidad de medida en base a su ID
func NewUM(id string) (UM, error) {
	// Chequeo que existan las unidades de medida
	d, ok := unidades[id]
	if !ok {
		return d, errors.Errorf("No existe la unidad de medida '%v'", id)
	}
	return d, nil
}

type RelacionDeUnidad interface {
	Factor(desde, hacia string) (factor float64)
}

// RelacionUM representa una equivalencia entre Unidades de distinto tipo.
// Por ejemplo: 1 bidon, equivale a 25 kg
//
// Si para un producto necesito convertir de litros a peso por ejemplo,
// el resultado va a depender del peso específico del líquido.
// Para ello le doy la opción a la función de conversión de que
// cada producto que solicite una conversión le agregue sus propias
// relaciones.
//
// Si las medidas que se quieren convertir no son del mismo tipo,
// Converitr() va a ir llamando a esta función hasta encontrar la
// conversión que satisfaga los tipos
type RelacionUM struct {
	Un        string  `json:"un"`
	EquivaleA float64 `json:"equivale_a"`
	De        string  `json:"de"`
}

// Convertir transforma (1784 gr => Kg) = 1,784
func Convertir(cantidad float64, desde string, hacia string, relaciones ...RelacionUM) (out float64, err error) {

	// Chequeo que existan las unidades de medida
	d, ok := unidades[desde]
	if !ok {
		return 0, errors.Errorf("No existe la unidad de medida %v", desde)
	}
	h, ok := unidades[hacia]
	if !ok {
		return 0, errors.Errorf("No existe la unidad de medida %v", hacia)
	}

	// Si son iguales no hago nada
	if desde == hacia {
		return cantidad, nil
	}

	if d.Tipo != h.Tipo {

		// Empiezo a analizar las relaciones ingresadas
		for _, v := range relaciones {

			// Quiero convertir 600ml a kg. La relación me dice que 1L = 0.92kg
			// Busco si la relación me convierte los tipos que estoy buscando.
			if unidades[v.Un].Tipo == d.Tipo && unidades[v.De].Tipo == h.Tipo {
				// Primero convierto 600 ml => L = 0.6
				nuevaUnidadConsistenteConDesde, err := Convertir(cantidad, desde, v.Un)
				if err != nil {
					return 0, errors.Wrap(err, "error interno derecho")
				}

				// Convierto 0.6 L => Kg
				nuevaUnidadFinal, err := Convertir(nuevaUnidadConsistenteConDesde, v.De, hacia)
				if err != nil {
					return 0, errors.Wrap(err, "error interno")
				}

				return nuevaUnidadFinal * v.EquivaleA, nil
			}

			// Si tengo la relación inversa:
			// Quiero convertir 600ml a kg. La relación me dice que 1kg = 0.92L
			if unidades[v.Un].Tipo == h.Tipo && unidades[v.De].Tipo == d.Tipo {
				relacionInvertida := RelacionUM{
					Un:        v.De,
					EquivaleA: 1 / v.EquivaleA,
					De:        v.Un,
				}
				return Convertir(cantidad, desde, hacia, relacionInvertida)

			}

		}

		// No eran convertibles y tampoco se ingresó una conversión
		return 0, errors.Errorf("No se puede convertir %v hacia %v. La primera es una medida de %v, la segunda en cambio es de %v.",
			d.Nombre, h.Nombre, d.Tipo, h.Tipo)
	}

	// Son convertibles
	return cantidad * d.factor / h.factor, nil
}

// TipoMedida representa un tipo de medida. Las diferentes unidades dentro de un tipo son
// convertibles por definición (gr a kg, cm a metro, etc).
type TipoMedida string

const (
	TipoCantidad   = "Cantidad"
	TipoPeso       = "Peso"
	TipoDistancia  = "Distancia"
	TipoSuperficie = "Superficie"
	TipoVolumen    = "Volumen"
)

// Validar devuelve true si es un tipo de medida aceptado.
func (u UM) Validar() error {

	res, ok := unidades[u.ID]
	if !ok {
		return errors.Errorf("La unidad de medida %v no está dentro de las permitidas.", u.ID)
	}
	if res.Tipo != u.Tipo {
		return errors.Errorf("La unidad %v es del tipo %v, no de %v", u.ID, res.Tipo, u.Tipo)
	}

	switch u.Tipo {
	case TipoCantidad, TipoPeso, TipoDistancia, TipoSuperficie, TipoVolumen:
		return nil
	default:
		return errors.Errorf("El tipo de medida %v no está dentro de los permitidos", u.Tipo)
	}
}

// MedidasDe devuelve las unidades de medida de "Peso" por ejemplo.
func MedidasDe(tipoMedida string) (medidas []UM, err error) {
	for _, v := range unidades {
		if v.Tipo == tipoMedida {
			medidas = append(medidas, v)
		}
	}
	return
}

// Medidas devuelve todas las medidas
func Medidas() (medidas map[string]UM) {
	medidas = map[string]UM{}
	for k, v := range unidades {
		medidas[k] = v
	}
	return medidas
}
