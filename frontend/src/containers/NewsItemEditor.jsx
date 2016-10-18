import moment from 'moment';
import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { createSelector, createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';
import * as utils from '../utils';

import Editor from './Editor';

const newsItemIdSelector = (state, props) => props.params.newsItemId;
const newsItemSelector = createSelector(
  selectors.newsItems,
  newsItemIdSelector,
  utils.get,
);

@connect(createStructuredSelector({
  newsItemId: newsItemIdSelector,
  newsItem: newsItemSelector,
}), actions)
export default class NewsItemEditor extends Component {
  static propTypes = {
    isNew: PropTypes.bool,
    target: PropTypes.string,
    targetId: PropTypes.string,
    newsItemId: PropTypes.string,
    newsItem: ImmutablePropTypes.map,
    onSuccess: PropTypes.func,
    noMessage: PropTypes.bool,

    createNewsItem: PropTypes.func.isRequired,
    loadNewsItem: PropTypes.func.isRequired,
    updateNewsItem: PropTypes.func.isRequired,
  };

  componentWillMount() {
    this.load(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.load(nextProps, this.props);
  }

  load({ isNew, newsItemId, newsItem }, prevProps) {
    if (isNew) return;
    if (
      !prevProps ||
      newsItemId !== prevProps.newsItemId ||
      (!newsItem && prevProps.newsItem)
    ) {
      this.props.loadNewsItem(newsItemId);
    }
  }

  createNewsItem = (body, ...args) => this.props.createNewsItem({
    target: this.props.target,
    targetId: this.props.targetId,
    ...body,
  }, ...args);

  render() {
    return (
      <Editor
        isNew={this.props.isNew}
        resourceName="news item"
        resource={this.props.newsItem}
        fields={{
          picture: { name: 'Picture', preview: 'image' },
          video: { name: 'YouTube video ID' },
          title: { name: 'Title', preview: 'title', standalone: true },
          body: {
            name: 'Body',
            type: 'textarea',
            preview: 'text',
            standalone: true,
          },
          publishedAt: {
            name: 'Publish time',
            type: 'datetime',
            value: moment((new Date()).setSeconds(0, 0)).toISOString(),
            onlyNew: true,
          },
        }}
        createResource={this.createNewsItem}
        updateResource={this.props.updateNewsItem}
        onSuccess={this.props.onSuccess}
        noMessage={this.props.noMessage}
      />
    );
  }
}
