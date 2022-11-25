import Route from '@ember/routing/route';
//import { service } from '@ember/service';
import { getOwner } from '@ember/application';

export default class RentalRoute extends Route {
  //@service store;
  get store() {
    return getOwner(this).lookup('service:store');
  }

  async model(params) {
    return this.store.findRecord('rental', params.rental_id);
  }
}
