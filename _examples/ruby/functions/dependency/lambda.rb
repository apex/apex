puts "start dependency function"

require 'active_support/core_ext/object/blank'

def handler(event:, context:)
  puts "event.blank?: #{event.blank?}"
end
