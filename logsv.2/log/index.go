package log

import (
	"io"
	"os"

	"github.com/tysonmote/gommap"
)

// el index es la parte donde se guardan índices (offsets) que permite tener las ubicaciones de la data dentro del store :0
var (
	offWidth uint64 = 4                   // es el tamaño del offset
	posWidth uint64 = 8                   //es el tamaño de la posición
	entWidth        = offWidth + posWidth //es el tamaño completo de lo que se guarda en el index
)

type index struct {
	file *os.File    // puntero al archivo
	mmap gommap.MMap // representará el archivo mapeado en memoria
	size uint64      // tamaño del archivo
}

func newIndex(f *os.File, c Config) (*index, error) {
	//sacamos info del archivo, es deecir, su tamaño
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}

	//guardamos el tamaño del archivo
	sizeFile := uint64(fi.Size())

	//ajusta tamaño del archivo  a su máxima capacidad
	if errTrunc := os.Truncate(f.Name(), int64(c.Segment.MaxIndexBytes)); errTrunc != nil {
		return nil, errTrunc
	}

	//mapeamos el archivo en memoria, como nos interesa leer y escribir, concedemos esos permisos, y permite el recurso compartido
	newmmap, errMap := gommap.Map(f.Fd(), gommap.PROT_READ|gommap.PROT_WRITE, gommap.MAP_SHARED)
	if errMap != nil {
		return nil, errMap
	}

	//devolvemos el índice
	return &index{
		file: f,
		mmap: newmmap,
		size: sizeFile,
	}, nil
}

// función que nos permite regresar el nombre del arhivo
func (i *index) Name() string {
	return i.file.Name()
}

// función para leer un registro guardado en el index, para eso nos pasan un offset
func (i *index) Read(idx int64) (out uint32, pos uint64, err error) {
	//ver si el offset no está vacio, si está vacío, error
	if i.size == 0 {
		return 0, 0, io.EOF
	}
	if idx == -1 {
		//si ingresa un idx -1, se obtiene el último offset dentro del index
		out = uint32((i.size / entWidth) - 1) //dividimos el archivo entre el tamaño de cada entrada para saber cuántas entradas hay y le restamos 1 (empezamos desde 0)
	} else {
		//si da un offset bien, solo se guarda ese offset
		out = uint32(idx)
	}

	//verificar que el tamaño del registro no esté fuera de rango
	if (uint64(out)*entWidth)+entWidth > i.size {
		return 0, 0, io.EOF
	} else {
		pos = uint64(out) * entWidth //nos ubicamos al comienzo del registro
	}
	//se leerá el offset del mapa de memoria, decodificado
	out = enc.Uint32(i.mmap[pos : pos+offWidth])
	//se leerá la posición del mapa de memoria , decodificado
	pos = enc.Uint64(i.mmap[pos+offWidth : pos+entWidth])

	return out, pos, nil

}

func (i *index) Write(off uint32, pos uint64) error {

	//verifiquemos si el nuevo dato cabe en memoria
	memo := uint64(len(i.mmap)) //obtenemos el tamaño de la memoria
	//si el tamaño de la memoria es menor a lo que hay ya en memoria + una nueva entrada, entonces no se puede
	if memo < i.size+entWidth {
		return io.EOF
	}

	//escribimos el offset en memoria
	enc.PutUint32(i.mmap[i.size:i.size+offWidth], off)
	//escribimos la posición después del offset en memoria
	enc.PutUint64(i.mmap[i.size+offWidth:i.size+entWidth], pos)
	//actualizamos el tamaño del index
	i.size += entWidth
	return nil
}

func (i *index) Close() error {

	//cerramos el archivo :D
	return i.file.Close()

}
