import { Track, TracksTableData } from '@type/track';
import { getProviderTextColor } from '@utils/getProviderStyle';

export const tracksToTableData = (tracks: Track[]): TracksTableData[] => {
  return tracks.map((track, index) => {
    return {
      id: track.id,
      content: [
        {
          type: "select",
          isSelected: track.isSelected,
        },
        {
          type: "index",
          index: index,
        },
        {
          type: "title",
          title: track.title,
          link: track.link,
          artist: track.artist.name,
          artistlink: track.artist.link,
          imageurl: track.picture,
        },
        {
          type: "text",
          text: track.platform,
          class: getProviderTextColor(track.platform),
          link: track.link,
        },
        {
          type: "text",
          text: track.album.title,
          link: track.album.link,
        },

        {
          type: "track",
          id: track.id,
          trackurl: track.preview,
        },
        {
          type: "duration",
          duration: track.duration_ms,
        },
      ],
    } as TracksTableData;
  });
};
