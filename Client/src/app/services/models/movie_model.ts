export interface Genre {
  genre_id: number;
  genre_name: string;
}

export interface Ranking {
  ranking_value: number;
  ranking_name: string;
}

export interface Movie {
  _id?: string;           // bson.ObjectID → string in frontend
  imdb_id: string;
  title: string;
  release_date: string;
  original_language: string;
  poster_path: string;
  youtube_id: string;
  genre: Genre[];
  admin_review?: string;
  ranking: Ranking;
}

export interface MovieResponse {
  movies: Movie[];
  page_number: number;
}

export interface GenresResponse {
  genres: Genre[];
}
