//import { Article as Article2 } from '../../components/article';
import { Article2 } from '../../components/article2';
import { loadPages } from './_load';

interface ElPaisArticle {
  byline: string;
  content: string[];
  excerpt: string;
  image: string;
  length: number;
  siteName: string;
  textContent: string;
  title: string;
  url: string;
}

interface ElPaisPage {
  id: string;
  'es-es': ElPaisArticle;
  'pt-br': ElPaisArticle;
}

export default function Page(props: { page: ElPaisPage }): any {
  const divStyle = {
    display: 'flex',
  };

  return (
    <div className="articles container mx-auto" style={divStyle}>
      <Article2
        article1={props.page['pt-br']}
        article2={props.page['es-es']}
      ></Article2>
    </div>
  );
}

export async function getStaticProps(context: any) {
  const res = loadPages();

  const page = res.find((a: any) => a.id == context.params.id);

  return {
    props: {
      page,
    }, // will be passed to the page component as props
  };
}

export async function getStaticPaths() {
  const res = loadPages();

  const params = res.map((a: any) => ({
    params: { id: a.id },
  }));

  return {
    paths: [...params],
    fallback: false,
  };
}
