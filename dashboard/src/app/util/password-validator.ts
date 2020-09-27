import { AbstractControl } from '@angular/forms';

export class PasswordValidator {
  static caseAlternatives(control: AbstractControl) {
    const regSmall = '[a-z]';
    const regCap = '[A-Z]';
    if (control.value.match(regSmall) && control.value.match(regCap)) {
      return null;
    } else {
      return { caseAlternatives: true };
    }
  }

  static haveNumbers(control: AbstractControl) {
    const reg = '[1-9]';
    if (control.value.match(reg)) {
      return null;
    } else {
      return { haveNumbers: true };
    }
  }

  static haveSymbols(control: AbstractControl) {
    const reg = '[!@#$%^&*.]';
    if (control.value.match(reg)) {
      return null;
    } else {
      return { haveSymbols: true };
    }
  }
}
