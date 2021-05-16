package elpais_test

import (
	"bilingual-articles/providers/elpais"
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRssGetter struct{}

func (m *mockRssGetter) Get(ctx context.Context, url string) (*gofeed.Feed, error) {
	f, err := os.Open("../../testdata/portaad-10-may-2021.xml")
	if err != nil {
		panic(err)
	}

	return gofeed.NewParser().Parse(f)
}

type mockHttpClient struct{}

func (m *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	f, err := os.Open("../../testdata/" + req.URL.Host + req.URL.Path)
	if err != nil {
		panic(err)
	}

	// TODO based on url
	return &http.Response{
		StatusCode: 200,
		Body:       f,
	}, nil
}

func TestElPaisFindBilingualPages(t *testing.T) {
	ep := elpais.Provider{
		HttpClient: &mockHttpClient{},
	}

	type args struct {
		url  string
		want *elpais.Page
	}

	now := time.Now()

	tc := []args{
		{
			url: "https://brasil.elpais.com/cultura/2021-05-09/agua-de-murta-o-desodorante-de-isabel-a-catolica.html",
			want: &elpais.Page{
				Published: &now,
				Provider:  "elpais", Links: []elpais.Link{
					{Url: "https://brasil.elpais.com/cultura/2021-05-09/agua-de-murta-o-desodorante-de-isabel-a-catolica.html", Lang: "pt-BR"},
					{Url: "https://elpais.com/cultura/2021-05-07/agua-de-murta-el-desodorante-de-isabel-la-catolica.html", Lang: "es-ES"}},
			},
		},

		{
			url: "https://brasil.elpais.com/internacional/2021-05-03/espanha-cobra-o-desbloqueio-do-acordo-da-uniao-europeia-com-o-mercosul.html",
			want: &elpais.Page{
				Published: &now,
				Provider:  "elpais", Links: []elpais.Link{{Url: "https://brasil.elpais.com/internacional/2021-05-03/espanha-cobra-o-desbloqueio-do-acordo-da-uniao-europeia-com-o-mercosul.html", Lang: "pt-BR"}, {Url: "https://elpais.com/internacional/2021-05-03/espana-reclama-a-bruselas-que-desbloquee-el-acuerdo-con-mercosur.html", Lang: "es-ES"}}},
		},

		// Article that does not have a bilingual version
		{
			url:  "http://brasil.elpais.com/brasil/2021/04/16/album/1618596565_470682.html",
			want: nil,
		},
	}

	for _, tt := range tc {
		t.Run(tt.url, func(t *testing.T) {
			got, err := ep.FindBilingualPages(
				context.TODO(),
				tt.url,
				&now,
			)
			assert.NoError(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestElPaisProcessPage(t *testing.T) {
	ep := elpais.Provider{
		HttpClient: &mockHttpClient{},
	}

	type args struct {
		name string
		page elpais.Page
		want elpais.ElPaisArticle
	}

	tc := []args{
		{
			page: elpais.Page{
				Provider: "elpais", Links: []elpais.Link{
					{Url: "https://brasil.elpais.com/cultura/2021-05-09/agua-de-murta-o-desodorante-de-isabel-a-catolica.html", Lang: "pt-BR"},
					{Url: "https://elpais.com/cultura/2021-05-07/agua-de-murta-el-desodorante-de-isabel-la-catolica.html", Lang: "es-ES"}},
			},
			want: elpais.ElPaisArticle{
				Id: "4297f3862a64394ed1367885671c604c",
				PtBr: elpais.ReadableArticle{
					Url:         "https://brasil.elpais.com/cultura/2021-05-09/agua-de-murta-o-desodorante-de-isabel-a-catolica.html",
					Title:       "Água de murta, o desodorante de Isabel, a Católica",
					Byline:      "Manuel Morales",
					TextContent: `Retrato de Isabel, a Católica, de autor anônimo, que está no Museu do Prado. O quadro foi feito por volta de 1490, época em que já tinha Sancho Paredes de Guzmán a seu serviçoSancho de Paredes Golfín foi um homem extremamente meticuloso, com um zelo por seu ofício que o fazia anotar tudo, um funcionário cuja obra foi um exercício de transparência na coroa espanhola. Sancho serviu como camareiro da rainha Isabel I de Castela de 1498 até novembro de 1504, quando esta faleceu. Seu trabalho permitia que ele acessasse as dependências da monarca, e anotou o que ela tinha para sua vida privada, e a mais íntima, em 10 livros de contas. “Essa prática já existiu em reis anteriores em Castela, mas o interessante é o que ele deixou por escrito”, diz por telefone Miguel Ángel Ladero Quesada, da Academia de História. “O camareiro era o chefe da casa e coordenava dezenas de pessoas sob seu comando”, acrescenta.Mais informaçõesO fio iniciado por Sancho de Paredes há mais de 500 anos não se rompeu. Suas anotações passaram de geração em geração até que a aristocrata Tatiana Pérez de Guzmán el Bueno decidiu, pouco antes de morrer em 2012, que seu imenso patrimônio, incluindo o arquivo, passaria a uma fundação com seu nome para que fosse conhecido. Esse fio se enovela agora no palácio dos Golfines de Abajo, no centro histórico de Cáceres, cidade no oeste da Espanha. Museu desde 2015, o palácio também abriga parte do patrimônio artístico da fundação. Os muros dessa casa solarenga foram levantados no século XV por um nobre, Alonso Golfín, o pai de Sancho, que ajudou os Reis Católicos em seu caminho à coroa.Lá, uma equipe de documentaristas “estuda, cataloga e digitaliza há nove anos cada folha do arquivo, do qual aproximadamente 2.400 fichas descritivas, compostas por mais de 90.000 páginas, podem ser consultadas no site da fundação”, diz a responsável pelo arquivo, Elisa Arroyo. Lá estão os 10 livros de contas de Sancho de Paredes. No nono, encadernado em pele, escreveu a relação de produtos para os perfumes usados por Isabel, a Católica, como a algália, “uma substância untuosa, de odor forte e sabor acre”, afirmou. E o almíscar, “um odor forte secretado pelo cervo almiscarado, utilizado como notas de fundo em perfumaria”. Para hidratar seu rosto, a rainha utilizava o benjoim, “uma resina de uma árvore das florestas tropicais do sudeste asiático”. Na penteadeira também havia “perfumes elaborados, como o âmbar fino, o óleo de azahar e a água de murta, utilizada como desodorante”. E um produto hoje muito usado, o óleo de rosa mosqueta, “para regenerar a pele e eliminar manchas, cicatrizes e estrias”.Página do livro de contas com a relação de alguns dos perfumes usados por Isabel, a Católica. ARCHIVO HISTÓRICO DE LA FUNDACIÓN TATIANA PÉREZ DE GUZMÁN EL BUENOEsses cosméticos mostram uma rainha preocupada com seu aspecto e higiene, uma imagem oposta à maledicência divulgada de que não gostava muito de se lavar. Ladero, professor de História Medieval, que fez sua tese sobre a conquista de Granada, frisa que essa crença de que Isabel, a Católica fez a promessa de não trocar de camisa até tomar o último bastião muçulmano na Espanha “é uma lenda urbana”. “Pelo contrário, sempre foi muito asseada. Primeiro por sua dignidade política, já que o rei era vigário de Deus e precisava se mostrar limpo. Também porque naquele ambiente político o asseio era símbolo de limpeza moral”.Nesse nono livro também se encontra a relação de brocados, veludos, tabuleiros de xadrez, instrumentos musicais, pinturas, peles (coelho, arminho e marta)... O segundo detalha as joias de ouro e prata (coroas, colares, correntes, braceletes, anéis...). No sexto, os toucados, mantéis, toalhas... Chapéus e sapatos no sétimo. O eficiente Sancho redigiu em outro volume um índice dos conteúdos de todos esses livros. Talvez pensasse que alguém precisaria dos conteúdos no futuro. “Tinha a vontade de passar à posteridade”, diz Arroyo. Por que pôde conservar documentos tão pessoais da rainha? “Quando Isabel morreu, ele entregou os livros à Controladoria de Contas e ficou com uma cópia ou talvez fosse o original, não se sabe”, diz o historiador.A documentação revela que havia uma relação próxima do funcionário com sua rainha. Está conservado o despacho em que os Reis Católicos ordenavam que Golfín e sua esposa, Isabel Cuello, também camareira real, recebessem “bom alojamento” e “a preços razoáveis” por onde passassem. Golfín e Cuello tiveram 16 filhos, “dos quais nove trabalharam na corte: pajens, escudeiros...”, indica a arquivista. E quando a Católica morreu, ele foi um dos testamentários que assinou o documento. O que aconteceu com todas essas posses da rainha anotadas por Sancho de Paredes? “A maioria foi vendida e utilizada para pagar dívidas pessoais de Isabel”, diz Ladero.Os papéis do arquivo, que foram apresentados recentemente na Academia de História, estão classificados pelas províncias pelas quais os Golfines passaram. Além de Cáceres, Córdoba, Ávila, Valencia, Granada, Madri, Salamanca... O trabalho de Sancho Paredes foi um monumento à burocracia. “É o arquivo da administração de uma família, os Golfines”, acrescenta Arroyo. “De modo que há morgados, capelanias, pleitos, testamentos, acordos matrimoniais, heranças, compras de edifícios, cartas de recomendação...”. Com o atrativo de que entre os séculos XIV e XVI os Golfines estiveram a serviço da monarquia, o que acrescenta a correspondência com reis. A documentação vai até boa parte do século XIX.Retrato de Sancho de Paredes Golfín na sala de armas do palácio dos Golfines de Abajo (Cáceres).ARCHIVO HISTÓRICO DE LA FUNDACIÓN TATIANA PÉREZ DE GUZMÁN EL BUENOO restante do museu dos Golfines é uma viagem veloz a cinco séculos de história da família: tapeçarias de Bruxelas do século XVII, luminárias de La Granja, baús de viagem do século XVIII, um salão de baile do século XIX... Mas a joia é a sala de armas, finalizada em 1509. Os Golfines, como se fossem reis, a decoraram com seus escudos heráldicos, pinturas murais e personagens de sua estirpe desenhados... e na parte alta, sob o artesoado policromado, gravaram uma inscrição: “ESTA OVRA MANDO FACER EL ONRADO CAVALLERO SANCHO DE PAREDES...” (Esta obra foi feita por ordem do honrado cavaleiro Sancho de Paredes), a orgulhosa ostentação do prestígio e poder desta antiga linhagem.`,
					Content:     []string{"<figure><img alt=\"Retrato de Isabel, a Católica, de autor anônimo, que está no Museu do Prado. O quadro foi feito por volta de 1490, época em que já tinha Sancho Paredes de Guzmán a seu serviço\" width=\"828\" height=\"1085\" loading=\"eager\" src=\"https://cloudfront-eu-central-1.images.arcpublishing.com/prisa/U5XJ4IPW5ZBU5NBIJUNFG5BKAE.jpeg\" srcset=\"\"/><figcaption><span>Retrato de Isabel, a Católica, de autor anônimo, que está no Museu do Prado. O quadro foi feito por volta de 1490, época em que já tinha Sancho Paredes de Guzmán a seu serviço</span><span></span></figcaption></figure>", "<p>Sancho de Paredes Golfín foi um homem extremamente meticuloso, com um zelo por seu ofício que o fazia anotar tudo, um funcionário cuja obra foi um exercício de transparência na <a href=\"https://brasil.elpais.com/noticias/espana/\" target=\"_blank\">coroa espanhola</a>. Sancho serviu como camareiro da <a href=\"https://brasil.elpais.com/noticias/monarquia/\" target=\"_blank\">rainha Isabel I de Castela</a> de 1498 até novembro de 1504, quando esta faleceu. Seu trabalho permitia que ele acessasse as dependências da monarca, e anotou o que ela tinha para sua vida privada, e a mais íntima, em 10 livros de contas. “Essa prática já existiu em reis anteriores em Castela, mas <a href=\"https://brasil.elpais.com/noticias/historia/\" target=\"_blank\">o interessante é o que ele deixou por escrito</a>”, diz por telefone Miguel Ángel Ladero Quesada, da Academia de História. “O camareiro era o chefe da casa e coordenava dezenas de pessoas sob seu comando”, acrescenta.</p>", "<section><h3>Mais informações</h3></section>", "<p>O fio iniciado por Sancho de Paredes há mais de 500 anos não se rompeu. Suas anotações passaram de geração em geração até que a aristocrata Tatiana Pérez de Guzmán el Bueno decidiu, pouco antes de morrer em 2012, que seu imenso patrimônio, incluindo o arquivo, passaria a uma fundação com seu nome para que fosse conhecido. Esse fio se enovela agora no palácio dos Golfines de Abajo, <a href=\"https://brasil.elpais.com/noticias/provincia-caceres/\" target=\"_blank\">no centro histórico de Cáceres, cidade no oeste da Espanha</a>. Museu desde 2015, o palácio também abriga parte do patrimônio artístico da fundação. Os muros dessa casa solarenga foram levantados no século XV por um nobre, Alonso Golfín, o pai de Sancho, que ajudou os Reis Católicos em seu caminho à coroa.</p>", "<p>Lá, uma equipe de documentaristas “estuda, cataloga e digitaliza há nove anos cada folha do arquivo, do qual aproximadamente 2.400 fichas descritivas, compostas por mais de 90.000 páginas, podem ser consultadas no site da fundação”, diz a responsável pelo arquivo, Elisa Arroyo. Lá estão os 10 livros de contas de Sancho de Paredes. No nono, encadernado em pele, escreveu a relação de produtos para os perfumes usados por Isabel, a Católica, como a algália, “uma substância untuosa, de odor forte e sabor acre”, afirmou.<b> </b>E o almíscar, “um odor forte secretado pelo cervo almiscarado, utilizado como notas de fundo em perfumaria”. Para hidratar seu rosto, a rainha utilizava o benjoim, “uma resina de uma árvore das florestas tropicais do sudeste asiático”. Na penteadeira também havia “perfumes elaborados, como o âmbar fino, o óleo de azahar e a água de murta, utilizada como desodorante”. E um produto hoje muito usado, o óleo de rosa mosqueta, “para regenerar a pele e eliminar manchas, cicatrizes e estrias”.</p>", "<figure><img alt=\"Página do livro de contas com a relação de alguns dos perfumes usados por Isabel, a Católica. \" width=\"2154\" height=\"3169\" loading=\"lazy\" src=\"https://imagens.brasil.elpais.com/resizer/ssbph9xOAGgfy5fBb5Gk1OR9xm8=/414x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/JRSZLSZQDNH4PNTATRKMDNUQBI.jpg\" srcset=\"https://imagens.brasil.elpais.com/resizer/ssbph9xOAGgfy5fBb5Gk1OR9xm8=/414x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/JRSZLSZQDNH4PNTATRKMDNUQBI.jpg 414w,https://imagens.brasil.elpais.com/resizer/4sQVWEBd8oenPeXYxLvHhyYTlRY=/828x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/JRSZLSZQDNH4PNTATRKMDNUQBI.jpg 640w,https://imagens.brasil.elpais.com/resizer/HvWxhAbYqOucqDj7AYS3EdSycPY=/980x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/JRSZLSZQDNH4PNTATRKMDNUQBI.jpg 1000w,https://imagens.brasil.elpais.com/resizer/RcIxGmYi315xBGpu5y3mE6AmY7c=/1960x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/JRSZLSZQDNH4PNTATRKMDNUQBI.jpg 1960w\"/><figcaption><span>Página do livro de contas com a relação de alguns dos perfumes usados por Isabel, a Católica. </span><span><span>ARCHIVO HISTÓRICO DE LA FUNDACIÓN TATIANA PÉREZ DE GUZMÁN EL BUENO</span></span></figcaption></figure>", "<p>Esses cosméticos mostram <a href=\"https://brasil.elpais.com/noticias/higiene/\" target=\"_blank\">uma rainha preocupada com seu aspecto e higiene</a>, uma imagem oposta à maledicência divulgada de que não gostava muito de se lavar. Ladero, professor de História Medieval, que fez sua tese sobre a conquista de Granada, frisa que essa crença de que Isabel, a Católica fez a promessa de não trocar de camisa até tomar o último bastião muçulmano na Espanha “é uma lenda urbana”. “Pelo contrário, sempre foi muito asseada. Primeiro por sua dignidade política, já que o rei era vigário de Deus e precisava se mostrar limpo. Também porque naquele ambiente político o asseio era símbolo de limpeza moral”.</p>", "<p>Nesse nono livro também se encontra a relação de brocados, veludos, tabuleiros de xadrez, instrumentos musicais, pinturas, peles (coelho, arminho e marta)... O segundo detalha as joias de ouro e prata (coroas, colares, correntes, braceletes, anéis...). No sexto, os toucados, mantéis, toalhas... Chapéus e sapatos no sétimo. O eficiente Sancho redigiu em outro volume um índice dos conteúdos de todos esses livros. Talvez pensasse que alguém precisaria dos conteúdos no futuro. “Tinha a vontade de passar à posteridade”, diz Arroyo. Por que pôde conservar documentos tão pessoais da rainha? “Quando Isabel morreu, ele entregou os livros à Controladoria de Contas e ficou com uma cópia ou talvez fosse o original, não se sabe”, diz o historiador.</p>", "<p>A documentação revela que havia uma relação próxima do funcionário com sua rainha. Está conservado o despacho em que os Reis Católicos ordenavam que Golfín e sua esposa, Isabel Cuello, também camareira real, recebessem “bom alojamento” e “a preços razoáveis” por onde passassem. Golfín e Cuello tiveram 16 filhos, “dos quais nove trabalharam na corte: pajens, escudeiros...”, indica a arquivista. E quando a Católica morreu, ele foi um dos testamentários que assinou o documento. O que aconteceu com todas essas posses da rainha anotadas por Sancho de Paredes? “A maioria foi vendida e utilizada para pagar dívidas pessoais de Isabel”, diz Ladero.</p>", "<p>Os papéis do arquivo, que foram apresentados recentemente na Academia de História, estão classificados pelas províncias pelas quais os Golfines passaram. Além de Cáceres, Córdoba, Ávila, Valencia, Granada, Madri, Salamanca... <a href=\"https://brasil.elpais.com/cultura/2021-05-08/peter-brown-pior-que-esquecer-a-historia-e-distorce-la-para-avivar-o-ressentimento.html\" target=\"_blank\">O trabalho de Sancho Paredes foi um monumento à burocracia</a>. “É o arquivo da administração de uma família, os Golfines”, acrescenta Arroyo. “De modo que há morgados, capelanias, pleitos, testamentos, acordos matrimoniais, heranças, compras de edifícios, cartas de recomendação...”. Com o atrativo de que entre os séculos XIV e XVI os Golfines estiveram a serviço da monarquia, o que acrescenta a correspondência com reis. A documentação vai até boa parte do século XIX.</p>", "<figure><img alt=\"Retrato de Sancho de Paredes Golfín na sala de armas do palácio dos Golfines de Abajo (Cáceres).\" width=\"5368\" height=\"2934\" loading=\"lazy\" src=\"https://imagens.brasil.elpais.com/resizer/umMzWcuKcQSPvSHSgH9sOy7uK0M=/414x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/URBHX74UC5B5PJ3GOWLB5JOU4I.jpg\" srcset=\"https://imagens.brasil.elpais.com/resizer/umMzWcuKcQSPvSHSgH9sOy7uK0M=/414x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/URBHX74UC5B5PJ3GOWLB5JOU4I.jpg 414w,https://imagens.brasil.elpais.com/resizer/tyTUpCJBm5ncgk_SCBGBCQRIPcI=/828x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/URBHX74UC5B5PJ3GOWLB5JOU4I.jpg 640w,https://imagens.brasil.elpais.com/resizer/I13QUUT47eMIMR_55JN8ZKBxZ_Q=/980x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/URBHX74UC5B5PJ3GOWLB5JOU4I.jpg 1000w,https://imagens.brasil.elpais.com/resizer/aLE557l5DPchkLYvxs9MSvRPhJU=/1960x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/URBHX74UC5B5PJ3GOWLB5JOU4I.jpg 1960w\"/><figcaption><span>Retrato de Sancho de Paredes Golfín na sala de armas do palácio dos Golfines de Abajo (Cáceres).</span><span><span>ARCHIVO HISTÓRICO DE LA FUNDACIÓN TATIANA PÉREZ DE GUZMÁN EL BUENO</span></span></figcaption></figure>", "<p>O restante do museu dos Golfines é uma viagem veloz a cinco séculos de história da família: tapeçarias de Bruxelas do século XVII, luminárias de La Granja, baús de viagem do século XVIII, um salão de baile do século XIX...<b> </b>Mas a joia é a sala de armas, finalizada em 1509. Os Golfines, como se fossem reis, a decoraram com seus escudos heráldicos, pinturas murais e personagens de sua estirpe desenhados... e na parte alta, sob o artesoado policromado, gravaram uma inscrição: “ESTA OVRA MANDO FACER EL ONRADO CAVALLERO SANCHO DE PAREDES...” (Esta obra foi feita por ordem do honrado cavaleiro Sancho de Paredes), a orgulhosa ostentação do prestígio e poder desta antiga linhagem.</p>"},
					Length:      6333,
					Excerpt:     "Arquivo de um funcionário da realeza, conservado por cinco séculos, revela a grande atenção que a rainha dava a sua higiene e aspecto, contrariando a lenda urbana que a tachou de desasseada",
					SiteName:    "EL PAÍS",
					Image:       "https://imagens.brasil.elpais.com/resizer/EFpQVmepnYuQpgSY1miWyNkBsiY=/1200x0/filters:focal(1777x2097:1787x2107)/cloudfront-eu-central-1.images.arcpublishing.com/prisa/U5XJ4IPW5ZBU5NBIJUNFG5BKAE.jpeg",
				},
				EsEs: elpais.ReadableArticle{
					Url:   "https://elpais.com/cultura/2021-05-07/agua-de-murta-el-desodorante-de-isabel-la-catolica.html",
					Title: "Agua de murta, el desodorante de Isabel la Católica", Byline: "Manuel Morales",
					Excerpt:     "El archivo de un funcionario de la realeza, conservado desde hace cinco siglos, desvela la gran atención que la soberana prestaba a su higiene y aspecto, en contra de la leyenda urbana que la tildó de desaseada",
					TextContent: `Retrato de Isabel la Católica, de autor anónimo, que está en el Museo del Prado. El cuadro fue realizado en torno a 1490, época en que ya tenía a su servicio a Sancho Paredes de Guzmán.Sancho de Paredes Golfín fue un hombre extremadamente meticuloso, con un celo por su oficio que le llevaba a apuntarlo todo, un funcionario cuya obra fue un ejercicio de transparencia en la corona española. Sancho sirvió como camarero de la reina Isabel I de Castilla desde 1498 hasta noviembre de 1504, cuando esta falleció. Su trabajo le permitía acceder a las dependencias de la monarca, de las que anotó lo que ella tenía para su vida privada, y la más íntima, en 10 libros de cuentas. “Ya hubo esa práctica en reyes anteriores en Castilla, pero lo interesante es que él lo dejó por escrito”, dice por teléfono Miguel Ángel Ladero Quesada, de la Academia de la Historia. “El camarero era el jefe de la casa y estaba al frente de decenas de personas que coordinaba”, añade.Más información'Una reina fausta e infausta'Las relaciones de España y Portugal a través del arte: de Isabel a IsabelAdiós a sus católicas majestadesEl hilo que tiró Sancho de Paredes hace más de 500 años no se rompió. Sus apuntes pasaron de generación en generación hasta que la aristócrata Tatiana Pérez de Guzmán el Bueno decidió, poco antes de morir en 2012, que su inmenso patrimonio, incluido el archivo, pasara a una fundación con su nombre para que se diera a conocer. Ese hilo se ovilla ahora en el palacio de los Golfines de Abajo, en el casco histórico de Cáceres. Museo desde 2015, el palacio atesora también parte del patrimonio artístico de la fundación. Los muros de esta casa solariega los levantó en el siglo XV un noble, Alonso Golfín, el padre de Sancho, que había ayudado a los Reyes Católicos en su camino a la corona.Allí, un equipo de documentalistas “estudia, cataloga y digitaliza desde hace nueve años cada hoja del archivo, del que unas 2.400 fichas descriptivas, compuestas por más de 90.000 páginas, pueden consultarse en la web de la fundación”, explica la responsable del archivo, Elisa Arroyo. Ahí están los 10 libros de cuentas de Sancho de Paredes. En el noveno, encuadernado en piel, escribió la relación de productos para perfumes que usaba Isabel la Católica, como la algalia, “una sustancia untuosa, de olor fuerte y sabor acre”, explicaba. O el almizcle, “un fuerte olor que segrega el macho del ciervo almizclero, utilizado como notas de fondo en perfumería”. Para hidratar su rostro, la reina recurría al benjuí, “una resina de un árbol de los bosques tropicales del sudeste asiático”. En el tocador había también “perfumes elaborados, como el ámbar fino, el aceite de azahar o el agua de murta, que se utilizaba como desodorante”. Y un producto hoy muy usado, el aceite de rosa mosqueta, “para regenerar la piel y eliminar manchas, cicatrices y estrías”.Página del libro de cuentas con la relación de algunos de los perfumes que usaba Isabel la Católica.ARCHIVO HISTÓRICO DE LA FUNDACIÓN TATIANA PÉREZ DE GUZMÁN EL BUENOEstos cosméticos muestran a una reina preocupada por su aspecto e higiene, una estampa opuesta a la maledicencia que difundió que no era muy aficionada a lavarse. Ladero, catedrático de Historia Medieval, que hizo su tesis sobre la conquista de Granada, subraya que esa creencia de que Isabel la Católica hizo la promesa de no cambiarse de camisa hasta tomar el último bastión musulmán en España “es una leyenda urbana”. “Al contrario, siempre fue muy aseada. Primero por su dignidad política, ya que el rey era vicario de Dios y tenía que mostrarse limpio. También porque en aquel ambiente político el aseo era símbolo de limpieza moral”.En el libro noveno también se incluye la relación de brocados, terciopelos, tableros de ajedrez, instrumentos musicales, pinturas, pieles (conejo, armiño y marta)... El segundo detalla las joyas de oro y plata (coronas, collares, cadenas, brazaletes, sortijas…). En el sexto, los tocados, manteles, toallas… Sombreros y zapatos en el séptimo. El eficiente Sancho redactó en otro volumen un índice de los contenidos de todos estos libros. Pensaba quizás en que alguien lo necesitase en el futuro. “Tenía el afán de pasar a la posteridad”, apunta Arroyo. ¿Por qué pudo conservar documentos tan personales de la reina? “Al morir Isabel, él entregó los libros a la Contaduría de Cuentas y se quedó con una copia o quizás era el original, no se sabe”, indica el historiador.Los Golfines estuvieron al servicio de la monarquía entre los siglos XIV y XVI De la documentación se desprende que había una relación cercana del funcionario con su señora. Se conserva la cédula en la que los Reyes Católicos ordenaban que a Golfín y a su esposa, Isabel Cuello, también camarera real, se les diese “buen alojamiento” y “a razonables precios” por donde pasasen. Golfín y Cuello tuvieron 16 hijos, “de los que nueve trabajaron en la corte: pajes, escuderos…”, indica la archivera. Y al expirar la Católica, él fue uno de los testamentarios que firmó el documento. ¿Qué pasó con todas esas posesiones de la reina que anotó Sancho de Paredes? “La mayoría se vendieron o se emplearon para pagar deudas personales de ella”, precisa Ladero.Los papeles del archivo, que se han presentado recientemente en la Academia de la Historia, están clasificados por las provincias donde los Golfines tuvieron presencia. Junto a Cáceres, Córdoba, Ávila, Valencia, Granada, Madrid, Salamanca… Lo de Sancho de Paredes fue un monumento a la burocracia. “Se trata del archivo de la administración de una familia, los Golfines”, añade Arroyo. “Así que hay mayorazgos, capellanías, pleitos, testamentos, capitulaciones matrimoniales, herencias, compras de edificios, cartas de recomendación…”. Con el aliciente de que entre los siglos XIV y XVI los Golfines estuvieron al servicio de la monarquía, lo que suma correspondencia con reyes. La documentación llega hasta bien entrado el XIX.Retrato de Sancho de Paredes Golfín en la sala de armas del palacio de los Golfines de Abajo (Cáceres).ARCHIVO HISTÓRICO DE LA FUNDACIÓN TATIANA PÉREZ DE GUZMÁN EL BUENOEl resto del museo de los Golfines es un viaje veloz a cinco siglos de historia de la familia: tapices de Bruselas del siglo XVII, lámparas de La Granja, baúles de viaje del XVIII, un decimonónico salón de baile… Pero la joya es la sala de armas, acabada en 1509. Los Golfines, como si fuesen reyes, la decoraron con sus escudos heráldicos, pinturas murales y personajes de su estirpe dibujados… Y en la parte alta, bajo el artesonado policromado, grabaron una inscripción: “ESTA OVRA MANDO FACER EL ONRADO CAVALLERO SANCHO DE PAREDES…”, la orgullosa ostentación del prestigio y poder de este antiguo linaje.Una fundación preocupada por el cerebroLa Fundación Tatiana Pérez de Guzmán el Bueno se dedica, sobre todo, “a apoyar la investigación científica, especialmente la neurociencia”, dice en su sede de Madrid su director académico, Álvaro Matud. “Financiamos laboratorios en toda España y a jóvenes predoctorales”. En su faceta humanística, organizan “ciclos de conferencias sobre la historia de España” y gracias a su taller de restauración dan a conocer obras del patrimonio artístico de la aristócrata donostiarra. A esto se suma una biblioteca que, solo en Madrid, alberga unos 5.000 volúmenes, más fotografías, muebles, objetos, vestidos…`,

					Content:  []string{"<figure><img alt=\"Retrato de Isabel la Católica, de autor anónimo, que está en el Museo del Prado. El cuadro fue realizado en torno a 1490, época en que ya tenía a su servicio a Sancho Paredes de Guzmán.\" width=\"828\" height=\"1085\" loading=\"eager\" src=\"https://cloudfront-eu-central-1.images.arcpublishing.com/prisa/U5XJ4IPW5ZBU5NBIJUNFG5BKAE.jpeg\" srcset=\"\"/><figcaption><span>Retrato de Isabel la Católica, de autor anónimo, que está en el Museo del Prado. El cuadro fue realizado en torno a 1490, época en que ya tenía a su servicio a Sancho Paredes de Guzmán.</span><span></span></figcaption></figure>", "<p><a href=\"http://dbe.rah.es/biografias/91072/sancho-de-paredes-golfin\" target=\"_blank\">Sancho de Paredes Golfín</a> fue un hombre extremadamente meticuloso, con un celo por su oficio que le llevaba a apuntarlo todo, un funcionario cuya obra fue un ejercicio de transparencia en la corona española. Sancho sirvió como camarero de la reina <a href=\"http://dbe.rah.es/biografias/13005/isabel-i\" target=\"_blank\">Isabel I de Castilla</a> desde 1498 hasta noviembre de 1504, cuando esta falleció. Su trabajo le permitía acceder a las dependencias de la monarca, de las que anotó lo que ella tenía para su vida privada, y la más íntima, en 10 libros de cuentas. “Ya hubo esa práctica en reyes anteriores en Castilla, pero lo interesante es que él lo dejó por escrito”, dice por teléfono <a href=\"https://www.ucm.es/amcytme-historiamedieval/ladero-quesada,-miguel-angel\" target=\"_blank\">Miguel Ángel Ladero Quesada</a>, de la Academia de la Historia. “El camarero era el jefe de la casa y estaba al frente de decenas de personas que coordinaba”, añade.</p>", "<section><h3>Más información</h3><ul><li><a href=\"https://elpais.com/diario/2004/11/03/opinion/1099436410_850215.html?rel=listapoyo\">&#39;Una reina fausta e infausta&#39;</a></li><li><a href=\"https://elpais.com/cultura/2014/10/22/actualidad/1413989734_360438.html?rel=listapoyo\">Las relaciones de España y Portugal a través del arte: de Isabel a Isabel</a></li><li><a href=\"https://elpais.com/cultura/2014/11/28/television/1417203309_426724.html?rel=listapoyo\">Adiós a sus católicas majestades</a></li></ul></section>", "<p>El hilo que tiró Sancho de Paredes hace más de 500 años no se rompió. Sus apuntes pasaron de generación en generación hasta que la aristócrata Tatiana Pérez de Guzmán el Bueno decidió, poco antes de morir en 2012, que su inmenso patrimonio, incluido el archivo, pasara a una <a href=\"https://fundaciontatianapgb.org/\" target=\"_blank\">fundación</a> con su nombre para que se diera a conocer. Ese hilo se ovilla ahora en el <a href=\"https://www.palaciogolfinesdeabajo.com/\" target=\"_blank\">palacio de los Golfines de Abajo</a>, en el casco histórico de Cáceres. Museo desde 2015, el palacio atesora también parte del patrimonio artístico de la fundación. Los muros de esta casa solariega los levantó en el siglo XV un noble, Alonso Golfín, el padre de Sancho, que había ayudado a los Reyes Católicos en su camino a la corona.</p>", "<p>Allí, un equipo de documentalistas “estudia, cataloga y digitaliza desde hace nueve años cada hoja del archivo, del que unas 2.400 fichas descriptivas, compuestas por más de 90.000 páginas, pueden consultarse <a href=\"https://archivohistorico.es/\" target=\"_blank\">en la web de la fundación”, </a>explica la responsable del archivo, Elisa Arroyo. Ahí están los 10 libros de cuentas de Sancho de Paredes. En el noveno, encuadernado en piel, escribió la relación de productos para perfumes que usaba Isabel la Católica, como la algalia, “una sustancia untuosa, de olor fuerte y sabor acre”, explicaba. O el almizcle, “un fuerte olor que segrega el macho del ciervo almizclero, utilizado como notas de fondo en perfumería”. Para hidratar su rostro, la reina recurría al benjuí, “una resina de un árbol de los bosques tropicales del sudeste asiático”. En el tocador había también “perfumes elaborados, como el ámbar fino, el aceite de azahar o el agua de murta, que se utilizaba como desodorante”. Y un producto hoy muy usado, el aceite de rosa mosqueta, “para regenerar la piel y eliminar manchas, cicatrices y estrías”.</p>", "<figure><img alt=\"Página del libro de cuentas con la relación de algunos de los perfumes que usaba Isabel la Católica.\" width=\"2154\" height=\"3169\" loading=\"lazy\" src=\"https://imagenes.elpais.com/resizer/ssbph9xOAGgfy5fBb5Gk1OR9xm8=/414x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/JRSZLSZQDNH4PNTATRKMDNUQBI.jpg\" srcset=\"https://imagenes.elpais.com/resizer/ssbph9xOAGgfy5fBb5Gk1OR9xm8=/414x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/JRSZLSZQDNH4PNTATRKMDNUQBI.jpg 414w,https://imagenes.elpais.com/resizer/4sQVWEBd8oenPeXYxLvHhyYTlRY=/828x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/JRSZLSZQDNH4PNTATRKMDNUQBI.jpg 640w,https://imagenes.elpais.com/resizer/HvWxhAbYqOucqDj7AYS3EdSycPY=/980x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/JRSZLSZQDNH4PNTATRKMDNUQBI.jpg 1000w,https://imagenes.elpais.com/resizer/RcIxGmYi315xBGpu5y3mE6AmY7c=/1960x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/JRSZLSZQDNH4PNTATRKMDNUQBI.jpg 1960w\"/><figcaption><span>Página del libro de cuentas con la relación de algunos de los perfumes que usaba Isabel la Católica.</span><span><span>ARCHIVO HISTÓRICO DE LA FUNDACIÓN TATIANA PÉREZ DE GUZMÁN EL BUENO</span></span></figcaption></figure>", "<p>Estos cosméticos muestran a una reina preocupada por su aspecto e higiene, una estampa opuesta a la maledicencia que difundió que no era muy aficionada a lavarse. Ladero, catedrático de Historia Medieval, que hizo su tesis sobre la conquista de Granada, subraya que esa creencia de que Isabel la Católica hizo la promesa de no cambiarse de camisa hasta tomar el último bastión musulmán en España “es una leyenda urbana”. “Al contrario, siempre fue muy aseada. Primero por su dignidad política, ya que el rey era vicario de Dios y tenía que mostrarse limpio. También porque en aquel ambiente político el aseo era símbolo de limpieza moral”.</p>", "<p>En el libro noveno también se incluye la relación de brocados, terciopelos, tableros de ajedrez, instrumentos musicales, pinturas, pieles (conejo, armiño y marta)... El segundo detalla las joyas de oro y plata (coronas, collares, cadenas, brazaletes, sortijas…). En el sexto, los tocados, manteles, toallas… Sombreros y zapatos en el séptimo. El eficiente Sancho redactó en otro volumen un índice de los contenidos de todos estos libros. Pensaba quizás en que alguien lo necesitase en el futuro. “Tenía el afán de pasar a la posteridad”, apunta Arroyo. ¿Por qué pudo conservar documentos tan personales de la reina? “Al morir Isabel, él entregó los libros a la Contaduría de Cuentas y se quedó con una copia o quizás era el original, no se sabe”, indica el historiador.</p>", "<blockquote><p>Los Golfines estuvieron al servicio de la monarquía entre los siglos XIV y XVI </p></blockquote>", "<p>De la documentación se desprende que había una relación cercana del funcionario con su señora. Se conserva la cédula en la que los Reyes Católicos ordenaban que a Golfín y a su esposa, Isabel Cuello, también camarera real, se les diese “buen alojamiento” y “a razonables precios” por donde pasasen. Golfín y Cuello tuvieron 16 hijos, “de los que nueve trabajaron en la corte: pajes, escuderos…”, indica la archivera. Y al expirar la Católica, él fue uno de los testamentarios que firmó el documento. ¿Qué pasó con todas esas posesiones de la reina que anotó Sancho de Paredes? “La mayoría se vendieron o se emplearon para pagar deudas personales de ella”, precisa Ladero.</p>", "<p>Los papeles del archivo, que se han presentado recientemente en la Academia de la Historia, están clasificados por las provincias donde los Golfines tuvieron presencia. Junto a Cáceres, Córdoba, Ávila, Valencia, Granada, Madrid, Salamanca… Lo de Sancho de Paredes fue un monumento a la burocracia. “Se trata del archivo de la administración de una familia, los Golfines”, añade Arroyo. “Así que hay mayorazgos, capellanías, pleitos, testamentos, capitulaciones matrimoniales, herencias, compras de edificios, cartas de recomendación…”. Con el aliciente de que entre los siglos XIV y XVI los Golfines estuvieron al servicio de la monarquía, lo que suma correspondencia con reyes. La documentación llega hasta bien entrado el XIX.</p>", "<figure><img alt=\"Retrato de Sancho de Paredes Golfín en la sala de armas del palacio de los Golfines de Abajo (Cáceres).\" width=\"5368\" height=\"2934\" loading=\"lazy\" src=\"https://imagenes.elpais.com/resizer/umMzWcuKcQSPvSHSgH9sOy7uK0M=/414x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/URBHX74UC5B5PJ3GOWLB5JOU4I.jpg\" srcset=\"https://imagenes.elpais.com/resizer/umMzWcuKcQSPvSHSgH9sOy7uK0M=/414x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/URBHX74UC5B5PJ3GOWLB5JOU4I.jpg 414w,https://imagenes.elpais.com/resizer/tyTUpCJBm5ncgk_SCBGBCQRIPcI=/828x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/URBHX74UC5B5PJ3GOWLB5JOU4I.jpg 640w,https://imagenes.elpais.com/resizer/I13QUUT47eMIMR_55JN8ZKBxZ_Q=/980x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/URBHX74UC5B5PJ3GOWLB5JOU4I.jpg 1000w,https://imagenes.elpais.com/resizer/aLE557l5DPchkLYvxs9MSvRPhJU=/1960x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/URBHX74UC5B5PJ3GOWLB5JOU4I.jpg 1960w\"/><figcaption><span>Retrato de Sancho de Paredes Golfín en la sala de armas del palacio de los Golfines de Abajo (Cáceres).</span><span><span>ARCHIVO HISTÓRICO DE LA FUNDACIÓN TATIANA PÉREZ DE GUZMÁN EL BUENO</span></span></figcaption></figure>", "<p>El resto del museo de los Golfines es un viaje veloz a cinco siglos de historia de la familia: tapices de Bruselas del siglo XVII, lámparas de La Granja, baúles de viaje del XVIII, un decimonónico salón de baile… Pero la joya es la sala de armas, acabada en 1509. Los Golfines, como si fuesen reyes, la decoraron con sus escudos heráldicos, pinturas murales y personajes de su estirpe dibujados… Y en la parte alta, bajo el artesonado policromado, grabaron una inscripción: “ESTA OVRA MANDO FACER EL ONRADO CAVALLERO SANCHO DE PAREDES…”, la orgullosa ostentación del prestigio y poder de este antiguo linaje.</p>", "<section><h3>Una fundación preocupada por el cerebro</h3><figure><img alt=\"Fachada del palacio de los Golfines de Abajo (Cáceres).\" width=\"4032\" height=\"3024\" loading=\"lazy\" src=\"https://imagenes.elpais.com/resizer/jaXeMiA5hDqrMNbTL9gqxhripUQ=/414x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/QEDKXXCWVVD6PBDLZQ7FC433ZU.jpg\" srcset=\"https://imagenes.elpais.com/resizer/jaXeMiA5hDqrMNbTL9gqxhripUQ=/414x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/QEDKXXCWVVD6PBDLZQ7FC433ZU.jpg 414w,https://imagenes.elpais.com/resizer/GLFOZ8LYhkhlkUkjWjut-Tk2o4w=/828x0/cloudfront-eu-central-1.images.arcpublishing.com/prisa/QEDKXXCWVVD6PBDLZQ7FC433ZU.jpg 640w\"/></figure><p>La Fundación Tatiana Pérez de Guzmán el Bueno se dedica, sobre todo, “a apoyar la investigación científica, especialmente la neurociencia”, dice en su sede de Madrid su director académico, Álvaro Matud. “Financiamos laboratorios en toda España y a jóvenes predoctorales”. En su faceta humanística, organizan “ciclos de conferencias sobre la historia de España” y gracias a su taller de restauración dan a conocer obras del patrimonio artístico de la aristócrata donostiarra. A esto se suma una biblioteca que, solo en Madrid, alberga unos 5.000 volúmenes, más fotografías, muebles, objetos, vestidos…</p></section>"},
					SiteName: "EL PAÍS",

					Length: 7325,
					Image:  "https://imagenes.elpais.com/resizer/EFpQVmepnYuQpgSY1miWyNkBsiY=/1200x0/filters:focal(1777x2097:1787x2107)/cloudfront-eu-central-1.images.arcpublishing.com/prisa/U5XJ4IPW5ZBU5NBIJUNFG5BKAE.jpeg"},
			},
		},
	}

	for _, tt := range tc {
		t.Run("a valid page", func(t *testing.T) {
			got, err := ep.ProcessPage(
				context.TODO(),
				tt.page,
			)
			assert.NoError(t, err)

			require.NotNil(t, tt.want.EsEs)
			require.NotNil(t, tt.want.PtBr)

			assert.Equal(t, &tt.want, got)

		})
	}
}