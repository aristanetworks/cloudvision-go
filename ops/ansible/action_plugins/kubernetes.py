import os
from ansible.plugins.action import ActionBase
from ansible.utils.vars import merge_hash

class ActionModule(ActionBase):

   def run(self, tmp=None, task_vars=None):
      self._supports_check_mode = True
      self._supports_async      = False

      result = super(ActionModule, self).run(tmp, task_vars)

      display = os.getenv('DISPLAY_K8S_YAML')
      if self._play_context.check_mode and display == self._task.name:
         # output the yaml command
         result['failed'] = True
         result['msg'] = "Stopping after yaml output"
         print "\x1B[35m================ YAML OUTPUT FOR TASK ================"
         print self._task.args['inline_data']
         print "======================================================\x1B[m"
      else:
         # Run the real kubernetes module
         result = merge_hash(
            result,
            self._execute_module(tmp=tmp, task_vars=task_vars),
         )

      return result
