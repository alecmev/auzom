import moment from 'moment';
import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import { createSelector, createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';
import * as utils from '../utils';
import Component from '../utils/Component';

import Comments from './Comments';

import Datetime from '../components/Datetime';
import Img from '../components/Img';
import Markdown from '../components/Markdown';

const newsItemIdSelector = (state, props) => props.params.newsItemId;
const newsItemSelector = createSelector(
  selectors.newsItems,
  newsItemIdSelector,
  utils.get,
);

@connect(createStructuredSelector({
  amAdmin: selectors.amAdmin,
  newsItemId: newsItemIdSelector,
  newsItem: newsItemSelector,
}), actions)
export default class SeasonNewsItem extends Component {
  static propTypes = {
    location: PropTypes.object, // eslint-disable-line

    amAdmin: PropTypes.bool,
    newsItemId: PropTypes.string,
    newsItem: ImmutablePropTypes.map,

    loadNewsItem: PropTypes.func.isRequired,
    patchNewsItem: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  state = {
    editing: false,
    playing: false,
  };

  componentWillMount() {
    this.load(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.load(nextProps, this.props);
  }

  load({ newsItemId, newsItem }, prevProps) {
    if (
      !prevProps ||
      newsItemId !== prevProps.newsItemId ||
      (!newsItem && prevProps.newsItem)
    ) {
      this.props.loadNewsItem(newsItemId);
    }
  }

  handleEdit = () => this.setState({ editing: true });

  handleEditSave = () => this.props.patchNewsItem(
    this.props.newsItemId, {
      action: 'publish',
      publishedAt: this.r.publishedAt.state.value,
    }, () => {
      this.setState({ editing: false });
      this.props.messagePush('updated publish time');
    },
  );

  handleEditDiscard = () => this.setState({ editing: false });

  handleDelete = () => {
    if (!this.props.newsItem) return;
    const { id, isDeleted } = this.props.newsItem.toJS();
    if (!isDeleted && !window.confirm('Delete this news article?')) return;
    this.props.patchNewsItem(id, { action: isDeleted ? 'undelete' : 'delete' });
  };

  handlePlay = () => this.setState({ playing: true });

  render() {
    const { amAdmin, newsItem } = this.props;
    if (!newsItem) return null;
    const { editing } = this.state;
    const {
      id, picture, title, video, body, isDeleted, ...x
    } = newsItem.toJS();
    const publishedAt = x.publishedAt && moment(x.publishedAt);
    const updatedAt = x.updatedAt && moment(x.updatedAt);
    return (
      <div className={this.cni({ u: 'sectionSingle' })}>
        {picture && !this.state.playing && <div className={this.cn({
          u: 'sectionMargined',
          d: 'pictureOuter',
        })}>
          <Img
            className={this.cn({ d: 'picture' })}
            src={picture}
          />
          {video && <button
            className={this.cn({ d: 'picturePlay' })}
            onClick={this.handlePlay}
          >
            <div className={this.cn({ d: 'picturePlayInner' })}>▶</div>
          </button>}
        </div>}
        {video && (this.state.playing || !picture) && <div className={this.cn({
          u: 'sectionMargined',
          d: 'videoOuter',
        })}>
          <iframe
            className={this.cn({ d: 'video' })}
            src={`//www.youtube.com/embed/${video}?` +
              'showinfo=0&' +
              'color=white&' +
              `autoplay=${picture ? 1 : 0}`
            }
            frameBorder="0"
            allowFullScreen
          />
        </div>}
        <div className={this.cn({ u: 'sectionMargined', d: 'inner' })}>
          <div className={this.cn({
            d: 'publishedAt',
            m: (!publishedAt || publishedAt.isAfter()) && 'hidden',
          })}>
            {publishedAt ? publishedAt.format('lll') : 'NO PUBLISH TIME'}
            {amAdmin && !isDeleted && !editing && <button
              className={this.cn({ d: 'metaAction' })}
              title="change publish time"
              onClick={this.handleEdit}
            >📅</button>}
            {amAdmin && editing && <button
              className={this.cn({ d: 'metaAction' })}
              title="save publish time"
              onClick={this.handleEditSave}
            >✔</button>}
            {amAdmin && editing && <button
              className={this.cn({ d: 'metaAction' })}
              title="discard publish time"
              onClick={this.handleEditDiscard}
            >✘</button>}
            {amAdmin && editing && <Datetime
              className={this.cn({ d: 'publishedAtInput' })}
              ref={this.rcb('publishedAt')}
              label="Publish time"
              value={publishedAt && publishedAt.toISOString()}
              open
            />}
          </div>
          <Markdown
            className={this.cn({ d: 'title' })}
            source={`# ${title}`}
          />
          <Markdown
            className={this.cn({ d: 'body' })}
            source={body}
          />
          {amAdmin && <div className={this.cn({ d: 'meta' })}>
            {updatedAt && <span
              className={this.cn({ d: 'metaEdited' })}
              title={updatedAt.format('lll')}
            >
              {isDeleted ? 'deleted' : 'last edited'} <b>
                {updatedAt.fromNow()}
              </b>
            </span>}
            {!isDeleted && <Link
              className={this.cn({ d: 'metaAction' })}
              to={`${this.props.location.pathname}/edit`}
              title="edit"
            >✎</Link>}
            <button
              className={this.cn({ d: 'metaAction' })}
              title={isDeleted ? 'undelete' : 'delete'}
              onClick={this.handleDelete}
            >{isDeleted ? '♻' : '×'}</button>
          </div>}
        </div>
        <Comments
          className={this.cn({ d: 'comments' })}
          target="news"
          targetId={id}
        />
      </div>
    );
  }
}
