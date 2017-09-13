import { PropTypes } from 'react';

import Component from '../utils/Component';

export default class Checkbox extends Component {
  static propTypes = {
    label: PropTypes.string,
    value: PropTypes.bool,
    onChange: PropTypes.func,
    onEnter: PropTypes.func,
    parentField: PropTypes.string, // stateless
    modified: PropTypes.bool,
  };

  constructor(props) {
    super(props);
    !this.props.parentField && (this.state = { value: !!props.value });
  }

  toggle = () => {
    const value = this.props.parentField ? !this.props.value : this.state.value;
    !this.props.parentField && this.setState({ value: !this.state.value });
    this.props.onChange && this.props.onChange(value, this.props.parentField);
  };

  handleClick = () => {
    this.toggle();
    this.r.me.blur();
  };

  handleKeyDown = (e) => {
    if (e.which === 32) e.preventDefault();
  };

  handleKeyUp = (e) => {
    if (e.which === 32) this.toggle();
    else if (this.props.onEnter && e.which === 13) this.props.onEnter();
  };

  render() {
    const value = this.props.parentField ? this.props.value : this.state.value;
    return (
      <div // eslint-disable-line
        ref={this.rcb('me')}
        className={this.cni()}
        tabIndex="0"
        onClick={this.handleClick}
        onKeyDown={this.handleKeyDown}
        onKeyUp={this.handleKeyUp}
      >
        <div className={this.cn({ d: 'inner' })}>
          <div className={this.cn({
            d: 'checkmark',
            s: value && 'checked',
          })} />
          <div className={this.cn({ d: 'label' })}>{this.props.label}</div>
        </div>
        <div
          className={this.cn({
            d: 'modified',
            s: [!this.props.modified && 'collapsed'],
          })}
          title="modified"
        />
      </div>
    );
  }
}
