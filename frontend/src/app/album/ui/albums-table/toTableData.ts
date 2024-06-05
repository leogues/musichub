import { Album, AlbumsTableData } from "app/album/album";

import { getProviderTextColor } from "@utils/getProviderStyle";

const link = ":provider/:albumId";

export const albumsToTableData = (albums: Album[]): AlbumsTableData[] => {
  return albums.map((album, index) => {
    return {
      id: album.id,
      isSelected: album.isSelected,
      link: link
        .replace(":provider", album.platform)
        .replace(":albumId", album.id),
      content: [
        {
          type: "index",
          index: index,
        },
        {
          type: "title",
          title: album.title,
          imageurl: album.picture,
        },
        {
          type: "text",
          text: album.platform,
          class: getProviderTextColor(album.platform),
        },
        {
          type: "text",
          text: album.artist.name,
        },

        {
          type: "text",
          text: `Faixas: ${album.total_tracks.toString()}`,
        },
        {
          type: "date",
          date: album.release_date,
        },
      ],
    } as AlbumsTableData;
  });
};
