import ImmutablePropTypes from 'react-immutable-proptypes';
import YouTube from 'react-youtube';
import scrollIntoViewIfNeeded from 'scroll-into-view-if-needed';

import Component from '../utils/Component';

export default class SeasonVideos extends Component {
  static propTypes = {
    season: ImmutablePropTypes.map,
  };

  state = {
    list: null,
    video: null,
  };

  componentWillMount() {
    this.load(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.load(nextProps, this.props);
  }

  load({ season }, prevProps) {
    if (!season || (prevProps && season === prevProps.season)) return;

    const listId = season.get('youtubePlaylist');
    if (!listId) return;

    fetch('https://www.googleapis.com/youtube/v3/playlistItems?' +
      `playlistId=${listId}&` +
      'part=snippet&' +
      'maxResults=11&' +
      'key=AIzaSyAcWyGmRE-2b7plIdTPU_9otkGtz2qNN_k',
    )
      .then(response => response.json())
      .then(json => this.setState({
        list: json.items.map(x => ({
          title: x.snippet.title,
          publishedAt: new Date(x.snippet.publishedAt),
          videoId: x.snippet.resourceId.videoId,
        })),
      }));
  }

  handleVideoSelectCache = {};
  handleVideoSelect = (video) => {
    if (!this.handleVideoSelectCache[video]) {
      this.handleVideoSelectCache[video] = () => this.setState({ video });
    }

    return this.handleVideoSelectCache[video];
  };

  playNext = () => {
    const { list, video } = this.state;
    const idx = list.findIndex(x => x.videoId === video);
    if (idx === -1 || idx === list.length - 1) return;
    this.setState({ video: list[idx + 1].videoId });
  };

  scrollActiveItemIntoView = x => x && scrollIntoViewIfNeeded(
    x, false, { duration: 150 },
  );

  render() {
    const { season } = this.props;
    if (!season) return null;
    const { list, video } = this.state;
    return (
      <div className={this.cni()}>
        <div className={this.cn({ d: 'inner' })}>
          <div className={this.cn({ d: 'videoOuter2' })}>
            <div className={this.cn({ d: 'videoOuter' })}>
              {video ? <YouTube
                className={this.cn({ d: 'video' })}
                videoId={video}
                opts={{
                  playerVars: {
                    showinfo: 0,
                    color: 'white',
                    autoplay: 1,
                    hd: 1,
                  },
                }}
                onEnd={this.playNext}
              /> : <div className={this.cn({ d: 'video' })}>
                NO VIDEO SELECTED
              </div>}
            </div>
          </div>
          <div className={this.cn({ d: 'list' })}>
            <div className={this.cn({ d: 'listInner' })}>
              {list && list.map(x => (
                <button
                  key={x.videoId}
                  ref={(x.videoId === video) && this.scrollActiveItemIntoView}
                  className={this.cn({
                    d: 'listItem',
                    s: x.videoId === video && 'selected',
                  })}
                  onClick={this.handleVideoSelect(x.videoId)}
                >
                  {x.title}
                </button>
              ))}
            </div>
          </div>
        </div>
      </div>
    );
  }
}
