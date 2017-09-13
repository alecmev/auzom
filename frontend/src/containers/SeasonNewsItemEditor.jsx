import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';

import NewsItemEditor from './NewsItemEditor';

export default class SeasonNewsItemEditor extends Component {
  static propTypes = {
    isNew: PropTypes.bool,
    season: ImmutablePropTypes.map,
    seasonPath: PropTypes.string,
    params: PropTypes.object.isRequired, // eslint-disable-line
  };

  static contextTypes = {
    router: PropTypes.object.isRequired,
  };

  handleSuccess = (x, y) => this.context.router.push(
    `${this.props.seasonPath}/news/${this.props.isNew ? x.id : y.id}`,
  );

  render() {
    const { season } = this.props;
    if (!season) return null;
    return (
      <NewsItemEditor
        isNew={this.props.isNew}
        target="season"
        targetId={season.get('id')}
        onSuccess={this.handleSuccess}
        params={this.props.params}
      />
    );
  }
}
